/*
 * This file is part of the go-demo application.
 * More license information can be found in the root folder.
 *
 * SPDX-License-Identifier: LicenseRef-TQSPSLA-1.0.3
 *
 * Copyright (c) 2025 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Maximilian Eschenbacher and the Energy Manager development team
 */

package modbus

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/goburrow/serial"
	"github.com/tbrandon/mbserver"
	serialClient "github.com/tq-systems/go-dbus/serial"
	"github.com/tq-systems/public-go-utils/v3/log"
	"github.com/tq-systems/public-go-utils/v3/mqtt"
)

type ModbusServer struct {
	serialClient     serialClient.Client
	server           *mbserver.Server
	mqttClient       mqtt.Client
	mqttSubscription mqtt.Subscription
	ticker           *time.Ticker
	done             chan struct{}
	registerMutex    sync.RWMutex
}

func NewModbusServer(mqttClient mqtt.Client, serialClient serialClient.Client) (*ModbusServer, error) {
	m := &ModbusServer{
		server:        mbserver.NewServer(),
		mqttClient:    mqttClient,
		serialClient:  serialClient,
		registerMutex: sync.RWMutex{},
	}

	// Initialize some demo values
	m.server.HoldingRegisters[0] = 42
	m.server.HoldingRegisters[1] = 0xdead
	m.server.HoldingRegisters[2] = 0xbeef

	// This register will be incremented by the ticker
	// It can be used to simulate a changing value in the Modbus server
	m.server.HoldingRegisters[3] = 0

	// Initialize coils
	m.server.Coils[0] = 1

	// Register custom Modbus function handler for function code 6 (Write Single Register)
	m.server.RegisterFunctionHandler(6, m.onWriteSingleRegister)

	// Set up Modbus TCP server
	if err := m.server.ListenTCP(DefaultModbusServerTcpAddrPort); err != nil {
		return nil, fmt.Errorf("modbus listen on %s: %w", DefaultModbusServerTcpAddrPort, err)
	}
	log.Debugf("Modbus server listening on %s", DefaultModbusServerTcpAddrPort)

	// If Modbus RTU is enabled, initialize it
	if EnableModbusServerRtu {
		err := m.initializeModbusRtu()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Modbus RTU: %v", err)
		}
		log.Debugf("Modbus server listening on /dev/tty%s", DefaultModbusServerRtuInterfaceName)
	}

	// Start a ticker that increments HoldingRegisters[3] every second
	m.done = make(chan struct{})
	m.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.registerMutex.Lock()
				m.server.HoldingRegisters[3]++
				m.registerMutex.Unlock()
			case <-m.done:
				return
			}
		}
	}()

	// Subscribe to MQTT topic for incoming messages
	subscription, err := mqttClient.Subscribe(MqttModbusServerInTopic, m.onMessage)
	if err != nil {
		return nil, fmt.Errorf("failed subscribing to mqtt topic %v: %v", MqttModbusServerInTopic, err)
	}
	m.mqttSubscription = subscription

	return m, nil
}

// initializeModbusRtu sets up the Modbus RTU interface and starts listening on it.
// It binds the serial interface and configures the Modbus server to listen on the specified serial port.
func (m *ModbusServer) initializeModbusRtu() error {
	err := m.serialClient.BindInterface(DefaultModbusServerRtuInterfaceName, nil, false)
	if err != nil {
		return fmt.Errorf("failed to bind serial interface %s: %v", DefaultModbusServerRtuInterfaceName, err)
	}

	err = m.server.ListenRTU(
		&serial.Config{
			Address:  "/dev/tty" + DefaultModbusServerRtuInterfaceName,
			BaudRate: DefaultModbusRtuBaudRate,
			DataBits: DefaultModbusRtuDataBits,
			StopBits: DefaultModbusRtuStopBits,
			Parity:   DefaultModbusRtuParity,
		})
	if err != nil {
		return fmt.Errorf("failed to listen on serial device: %v", err)
	}

	return nil
}

// onMessage handles incoming MQTT messages.
// It unmarshals the message, updates the corresponding Modbus register.
func (m *ModbusServer) onMessage(topic string, data []byte) {
	var message modbusServerMqttMessage
	if err := json.Unmarshal(data, &message); err != nil {
		log.Error("cannot unmarshal message: ", err)
		return
	}

	m.registerMutex.Lock()
	if int(message.Addr) < len(m.server.HoldingRegisters) {
		m.server.HoldingRegisters[message.Addr] = message.Value
	}
	m.registerMutex.Unlock()
}

// onWriteSingleRegister handles Modbus function code 6 (Write Single Register).
// It updates the HoldingRegisters with the provided address and value from the request.
// It also publishes the updated value to the MQTT topic.
func (m *ModbusServer) onWriteSingleRegister(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	// Modbus function 6 request PDU (excluding function code):
	// [0..1]=address (BE), [2..3]=value (BE)
	data := frame.GetData()
	log.Debugf("Received Modbus function %d request: %+v", frame.GetFunction(), frame)
	if len(data) != 4 {
		return nil, &mbserver.IllegalDataValue
	}

	addr := binary.BigEndian.Uint16(data[0:2])
	val := binary.BigEndian.Uint16(data[2:4])

	if int(addr) >= len(s.HoldingRegisters) {
		return nil, &mbserver.IllegalDataAddress
	}

	m.registerMutex.Lock()
	m.server.HoldingRegisters[addr] = val
	m.registerMutex.Unlock()
	m.publishMqttMsg(addr, val)

	// Per spec, response for function 6 echoes address+value
	return data[:4], &mbserver.Success
}

// publishMqttMsg publishes a message to the MQTT topic with the updated Modbus register value.
func (m *ModbusServer) publishMqttMsg(addr uint16, value uint16) {
	msg := modbusServerMqttMessage{
		Addr:  addr,
		Value: value,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error("cannot marshal MQTT message: ", err)
		return
	}
	err = m.mqttClient.PublishRaw(MqttModbusServerOutTopic, 2, false, data)
	if err != nil {
		log.Error("cannot publish MQTT message: ", err)
	}
}

// Destructor cleans up the Modbus server resources.
// It attempts a graceful shutdown first, and if that fails, it forcefully tears down the server.
// It also stops the ticker and closes the MQTT subscription.
// If Modbus RTU is enabled, it unbinds the serial interface.
// This method should be called when the Modbus server is no longer needed.
func (m *ModbusServer) Destructor() {
	m.mqttSubscription.Unsubscribe()
	if m.ticker != nil {
		m.ticker.Stop()
	}
	if m.done != nil {
		close(m.done)
	}
	// Attempt clean close first with timeout; if it blocks, perform brutal close.
	// Reason: mbserver.Close() can block indefinitely if RTU with no timeout is used.
	// There is a pull request for this issue in the project (tbrandon/mbserver) that has not yet been accepted.
	if err := closeServerWithTimeout(m.server, 500*time.Millisecond); err != nil {
		brutalCloseServer(m.server)
	}
	if EnableModbusServerRtu {
		err := m.serialClient.UnbindInterface(DefaultModbusServerRtuInterfaceName)
		if err != nil {
			log.Error("failed to unbind serial interface: ", err)
		}
	}
}

// closeServerWithTimeout attempts to gracefully close the Modbus server within a specified timeout.
// If the server does not close within the timeout, it returns an error.
func closeServerWithTimeout(s *mbserver.Server, d time.Duration) error {
	done := make(chan struct{})
	go func() {
		s.Close()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(d):
		return fmt.Errorf("mbserver.Close timeout after %s", d)
	}
}

// brutalCloseServer forcefully tears down mbserver.Server internals to avoid indefinite block.
// LAST RESORT: uses unsafe reflection; only called if normal Close() times out.
func brutalCloseServer(s *mbserver.Server) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("brutalCloseServer recovered: %v", r)
		}
	}()
	if s == nil {
		return
	}
	rv := reflect.ValueOf(s).Elem()
	// Close listeners
	if lf := rv.FieldByName("listeners"); lf.IsValid() && lf.Kind() == reflect.Slice {
		for i := 0; i < lf.Len(); i++ {
			lVal := lf.Index(i)
			if lVal.IsNil() {
				continue
			}
			if closer, ok := lVal.Interface().(interface{ Close() error }); ok {
				_ = closer.Close()
			}
		}
	}
	// Close ports directly
	if pf := rv.FieldByName("ports"); pf.IsValid() && pf.Kind() == reflect.Slice {
		for i := 0; i < pf.Len(); i++ {
			pVal := pf.Index(i)
			if pVal.IsNil() {
				continue
			}
			if closer, ok := pVal.Interface().(interface{ Close() error }); ok {
				_ = closer.Close()
			}
		}
	}
	// Close portsCloseChan if open
	if cf := rv.FieldByName("portsCloseChan"); cf.IsValid() {
		chPtr := unsafe.Pointer(cf.UnsafeAddr())
		chVal := *(*chan struct{})(chPtr)
		select {
		case <-chVal:
		default:
			// Try to close (will panic if already closed; protect with recover outer defer)
			close(chVal)
		}
	}
	// Drain requestChan non-blocking to allow handler to exit
	if rf := rv.FieldByName("requestChan"); rf.IsValid() {
		rcPtr := unsafe.Pointer(rf.UnsafeAddr())
		rc := *(*chan *mbserver.Request)(rcPtr)
		// Close requestChan to unblock handler goroutine if it's waiting
		// Protect with recover
		func() {
			defer func() { _ = recover() }()
			close(rc)
		}()
	}
	// Yield to let goroutines observe closes
	runtime.Gosched()
	// Attempt to wait a short moment
	time.Sleep(20 * time.Millisecond)
	// Attempt to zero WaitGroup by brute force not done (risk > benefit) – rely on GC after routines exit.
	log.Debugf("brutalCloseServer executed")
}
