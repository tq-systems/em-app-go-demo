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

import "time"

const (
	// Default Modbus RTU parameters
	DefaultModbusRtuBaudRate = 19200 // Baud rate for Modbus RTU
	DefaultModbusRtuDataBits = 8     // Data bits for Modbus RTU
	DefaultModbusRtuStopBits = 1     // Stop bits for Modbus RTU
	DefaultModbusRtuParity   = "N"   // Parity for Modbus RTU (N=none, E=even, O=odd)

	// Modbus Client settings
	DefaultModbusClientTimeout = 3 * time.Second // Default timeout for a single Modbus transaction

	// Modbus Server settings
	DefaultModbusServerTcpAddrPort      = "0.0.0.0:502" // Default TCP address and port for Modbus server
	DefaultModbusServerRtuInterfaceName = "APP1"         // Name of the serial interface for Modbus RTU server
	EnableModbusServerRtu               = true           // Set to true to enable Modbus RTU over serial interface

	// MQTT topics
	MqttModbusClientInTopic  = "modbus/local/values/client/in"
	MqttModbusClientOutTopic = "modbus/local/values/client/out"
	MqttModbusServerInTopic  = "modbus/local/values/server/in"
	MqttModbusServerOutTopic = "modbus/local/values/server/out"
)

// modbusServerMqttMessage represents the structure of the MQTT message sent/received by the Modbus server
// It contains the address and value of the Modbus register.
// The address is a 16-bit unsigned integer, and the value is also a 16-bit unsigned integer.
// This structure is used for both sending and receiving messages via MQTT.
type modbusServerMqttMessage struct {
	Addr  uint16 `json:"addr"`
	Value uint16 `json:"value"`
}

// Incoming MQTT request schema
// Example (TCP): {"protocol":"tcp","address":"192.168.1.10:502","unitId":1,"start":0,"quantity":4,"function":"holding"}
// Example (RTU): {"protocol":"rtu","device":"/dev/ttyAPP1","baud":19200,"dataBits":8,"stopBits":1,"parity":"N","unitId":1,"start":0,"quantity":2,"function":"input"}
// function: one of "coil","discrete","holding","input" mapping to FC 1,2,3,4
type modbusClientRequest struct {
	Protocol  string `json:"protocol"`          // "tcp" | "rtu"
	Address   string `json:"address,omitempty"` // host:port for TCP
	Device    string `json:"device,omitempty"`  // serial device for RTU
	Baud      int    `json:"baud,omitempty"`
	DataBits  int    `json:"dataBits,omitempty"`
	StopBits  int    `json:"stopBits,omitempty"`
	Parity    string `json:"parity,omitempty"`
	UnitID    byte   `json:"unitId"`
	Start     uint16 `json:"start"`
	Quantity  uint16 `json:"quantity"`
	Function  string `json:"function"` // coil|discrete|holding|input
	TimeoutMs int    `json:"timeoutMs,omitempty"`
	// Write support
	Op     string   `json:"op,omitempty"`     // "read" (default if omitted) | "write"
	Values []uint16 `json:"values,omitempty"` // for holding register writes (single: len=1, multi: len>1)
	Bits   []bool   `json:"bits,omitempty"`   // for coil writes (single: len=1, multi: len>1)
}

// Outgoing MQTT response schema
type modbusClientResponse struct {
	Ok         bool     `json:"ok"`
	Error      string   `json:"error,omitempty"`
	Protocol   string   `json:"protocol"`
	Address    string   `json:"address,omitempty"`
	Device     string   `json:"device,omitempty"`
	UnitID     byte     `json:"unitId"`
	Function   string   `json:"function"`
	Start      uint16   `json:"start"`
	Quantity   uint16   `json:"quantity"`
	Data       []uint16 `json:"data,omitempty"` // for holding/input reads
	Bits       []bool   `json:"bits,omitempty"` // for coil/discrete reads
	DurationMs int64    `json:"durationMs"`
	Op         string   `json:"op,omitempty"`
}
