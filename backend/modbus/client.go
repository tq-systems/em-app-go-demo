/*
 * Copyright (c) 2025-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229
 * Seefeld, Germany. All rights reserved.
 * Author: Maximilian Eschenbacher and the Energy Manager development team
 *
 * This software is licensed under the TQ-Systems Product Software License
 * Agreement Version 1.0.3 or any later version.
 * You can obtain a copy of the License Agreement in the TQS (TQ-Systems
 * Software Licenses) folder on the following website:
 * https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
 * In case of any license issues please contact license@tq-group.com.
 */

package modbus

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/goburrow/modbus"
	serialClient "github.com/tq-systems/go-dbus/serial"
	"github.com/tq-systems/public-go-utils/v3/log"
	"github.com/tq-systems/public-go-utils/v3/mqtt"
)

const (
	modbusReadOperation  = "read"
	modbusWriteOperation = "write"
)

type ModbusClient struct {
	serialClient     serialClient.Client
	mqttClient       mqtt.Client
	mqttSubscription mqtt.Subscription
	reqQueue         chan []byte
	quit             chan struct{}
}

func NewModbusClient(mqttClient mqtt.Client, serialClient serialClient.Client) (*ModbusClient, error) {
	m := &ModbusClient{
		serialClient: serialClient,
		mqttClient:   mqttClient,
		reqQueue:     make(chan []byte, 64),
		quit:         make(chan struct{}),
	}
	sub, err := mqttClient.Subscribe(MqttModbusClientInTopic, m.onMessage)
	if err != nil {
		return nil, fmt.Errorf("failed subscribing to mqtt topic %s: %v", MqttModbusClientInTopic, err)
	}
	m.mqttSubscription = sub
	// Start worker
	go m.worker()
	return m, nil
}

// Destructor unsubscribes MQTT and stops the worker goroutine.
func (m *ModbusClient) Destructor() {
	if m.mqttSubscription != nil {
		m.mqttSubscription.Unsubscribe()
	}
	close(m.quit)
	close(m.reqQueue)
}

// onMessage enqueues payload quickly to avoid doing heavy work in MQTT callback.
func (m *ModbusClient) onMessage(topic string, payload []byte) {
	select {
	case m.reqQueue <- payload:
	default:
		// Queue full -> drop oldest policy: discard this payload
		log.Warningf("Modbus client request queue full; dropping message")
	}
}

func (m *ModbusClient) worker() {
	for {
		select {
		case <-m.quit:
			return
		case payload, ok := <-m.reqQueue:
			if !ok {
				return
			}
			m.processPayload(payload)
		}
	}
}

func (m *ModbusClient) processPayload(payload []byte) {
	req, err := parseRequest(payload)
	if err != nil {
		m.publishError("", "", 0, 0, 0, err.Error())
		return
	}
	resp := initResponse(&req)
	if err := validateRequest(&req, resp.Op == modbusWriteOperation); err != nil {
		m.publishError(req.Protocol, req.Address, req.Start, req.Quantity, req.UnitID, err.Error())
		return
	}
	handler, cleanup, err := m.createHandler(&req)
	if err != nil {
		m.publishError(req.Protocol, req.Address, req.Start, req.Quantity, req.UnitID, err.Error())
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	client, setTimeoutOk := buildClient(handler)
	if client == nil {
		m.publishError(req.Protocol, req.Address, req.Start, req.Quantity, req.UnitID, "internal: handler type assertion failed")
		return
	}
	configureTimeout(handler, req.TimeoutMs)
	start := time.Now()
	var opErr error
	if resp.Op == modbusWriteOperation {
		opErr = executeWrite(client, &req, resp)
	} else {
		opErr = executeRead(client, &req, resp)
	}
	finalizeResponse(resp, opErr, start)
	_ = setTimeoutOk // placeholder if needed later
	m.publishResponse(resp)
}

// --- Helper functions for processPayload ---

func parseRequest(payload []byte) (modbusClientRequest, error) {
	var req modbusClientRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return req, fmt.Errorf("invalid JSON: %v", err)
	}
	return req, nil
}

func initResponse(req *modbusClientRequest) *modbusClientResponse {
	return &modbusClientResponse{
		Protocol: req.Protocol,
		Address:  req.Address,
		Device:   req.Device,
		UnitID:   req.UnitID,
		Function: req.Function,
		Start:    req.Start,
		Quantity: req.Quantity,
		Op:       chooseString(req.Op, modbusReadOperation),
	}
}

func validateRequest(req *modbusClientRequest, isWrite bool) error {
	if !isWrite && req.Quantity == 0 {
		return fmt.Errorf("quantity must be > 0")
	}
	if isWrite {
		switch req.Function {
		case "holding":
			if len(req.Values) == 0 {
				return fmt.Errorf("values required for holding write")
			}
		case "coil":
			if len(req.Bits) == 0 {
				return fmt.Errorf("bits required for coil write")
			}
		default:
			return fmt.Errorf("write not supported for function")
		}
	}
	return nil
}

func buildClient(handler interface{}) (modbus.Client, bool) {
	clientHandler, ok := handler.(modbus.ClientHandler)
	if !ok {
		return nil, false
	}
	return modbus.NewClient(clientHandler), true
}

func configureTimeout(handler interface{}, timeoutMs int) {
	deadline := DefaultModbusClientTimeout
	if timeoutMs > 0 {
		deadline = time.Duration(timeoutMs) * time.Millisecond
	}
	switch h := handler.(type) {
	case *modbus.TCPClientHandler:
		h.Timeout = deadline
	case *modbus.RTUClientHandler:
		h.Timeout = deadline
	}
}

func executeRead(client modbus.Client, req *modbusClientRequest, resp *modbusClientResponse) error {
	switch req.Function {
	case "holding":
		raw, err := client.ReadHoldingRegisters(req.Start, req.Quantity)
		if err != nil {
			return err
		}
		resp.Data = bytesToUint16(raw)
	case "input":
		raw, err := client.ReadInputRegisters(req.Start, req.Quantity)
		if err != nil {
			return err
		}
		resp.Data = bytesToUint16(raw)
	case "coil":
		raw, err := client.ReadCoils(req.Start, req.Quantity)
		if err != nil {
			return err
		}
		resp.Bits = bytesToBools(raw, int(req.Quantity))
	case "discrete":
		raw, err := client.ReadDiscreteInputs(req.Start, req.Quantity)
		if err != nil {
			return err
		}
		resp.Bits = bytesToBools(raw, int(req.Quantity))
	default:
		return fmt.Errorf("unknown function")
	}
	return nil
}

func executeWrite(client modbus.Client, req *modbusClientRequest, resp *modbusClientResponse) error {
	switch req.Function {
	case "holding":
		if len(req.Values) == 1 {
			_, err := client.WriteSingleRegister(req.Start, req.Values[0])
			if err != nil {
				return err
			}
		} else {
			l := len(req.Values)
			if l > math.MaxUint16 {
				return fmt.Errorf("values length exceeds uint16 range")
			}
			if req.Quantity != 0 && req.Quantity != uint16(l) {
				return fmt.Errorf("quantity mismatch values length")
			}
			req.Quantity = uint16(l)
			b := uint16sToBytes(req.Values)
			if _, err := client.WriteMultipleRegisters(req.Start, req.Quantity, b); err != nil {
				return err
			}
			resp.Quantity = req.Quantity
		}
		resp.Data = req.Values
	case "coil":
		if len(req.Bits) == 1 {
			var v uint16
			if req.Bits[0] {
				v = 0xFF00
			}
			if _, err := client.WriteSingleCoil(req.Start, v); err != nil {
				return err
			}
		} else {
			l := len(req.Bits)
			if l > math.MaxUint16 {
				return fmt.Errorf("bits length exceeds uint16 range")
			}
			if req.Quantity != 0 && req.Quantity != uint16(l) {
				return fmt.Errorf("quantity mismatch bits length")
			}
			packed := boolsToBytes(req.Bits)
			if _, err := client.WriteMultipleCoils(req.Start, uint16(l), packed); err != nil {
				return err
			}
			resp.Quantity = uint16(l)
		}
		resp.Bits = req.Bits
	default:
		return fmt.Errorf("write not implemented for function")
	}
	return nil
}

func finalizeResponse(resp *modbusClientResponse, opErr error, start time.Time) {
	if opErr != nil {
		resp.Ok = false
		resp.Error = opErr.Error()
	} else {
		resp.Ok = true
	}
	resp.DurationMs = time.Since(start).Milliseconds()
}

func (m *ModbusClient) createHandler(req *modbusClientRequest) (interface{}, func(), error) {
	switch req.Protocol {
	case "tcp":
		if req.Address == "" {
			return nil, nil, fmt.Errorf("address required for tcp")
		}
		h := modbus.NewTCPClientHandler(req.Address)
		h.SlaveId = req.UnitID
		if err := h.Connect(); err != nil {
			return nil, nil, fmt.Errorf("tcp connect: %v", err)
		}
		cleanup := func() { h.Close() }
		return h, cleanup, nil
	case "rtu":
		if req.Device == "" {
			return nil, nil, fmt.Errorf("device required for rtu")
		}
		// Bind interface if serialClient present
		if m.serialClient != nil {
			if err := m.serialClient.BindInterface(req.Device, nil, false); err != nil {
				return nil, nil, fmt.Errorf("bind interface %s failed: %v", req.Device, err)
			}
		}
		h := modbus.NewRTUClientHandler(req.Device)
		h.BaudRate = chooseInt(req.Baud, DefaultModbusRtuBaudRate)
		h.DataBits = chooseInt(req.DataBits, DefaultModbusRtuDataBits)
		h.Parity = chooseString(req.Parity, DefaultModbusRtuParity)
		h.StopBits = chooseInt(req.StopBits, DefaultModbusRtuStopBits)
		h.SlaveId = req.UnitID
		h.Timeout = DefaultModbusClientTimeout
		if err := h.Connect(); err != nil {
			return nil, nil, fmt.Errorf("rtu connect: %v", err)
		}
		cleanup := func() {
			h.Close()
			if m.serialClient != nil {
				_ = m.serialClient.UnbindInterface(req.Device)
			}
		}
		return h, cleanup, nil
	default:
		return nil, nil, fmt.Errorf("unknown protocol: %s", req.Protocol)
	}
}

func (m *ModbusClient) publishError(protocol, address string, start, quantity uint16, unitId byte, msg string) {
	resp := modbusClientResponse{Ok: false, Error: msg, Protocol: protocol, Address: address, Start: start, Quantity: quantity, UnitID: unitId}
	m.publishResponse(&resp)
}

func (m *ModbusClient) publishResponse(resp *modbusClientResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Error("cannot marshal Modbus client response: ", err)
		return
	}
	if err = m.mqttClient.PublishRaw(MqttModbusClientOutTopic, 1, false, data); err != nil {
		log.Error("cannot publish Modbus client response: ", err)
	}
}
