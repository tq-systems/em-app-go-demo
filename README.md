# Go Demo Application

> **Disclaimer (Demo Only)**: This code is provided **solely for demonstration and educational purposes** within the TQ Energy Manager SDK. It is **not production‑ready**. Before using it in a production system you must review, harden, test, and adapt it (error handling, security, resource management, performance, observability, configuration, protocol edge cases).

## Description
This is a comprehensive Go application for demonstration purposes designed for the TQ-Systems Energy Manager. It showcases best practices for developing Go applications on the Energy Manager platform, including MQTT communication, REST API endpoints, a frontend that displays GDR data from the backend in tabular form, Modbus TCP/RTU client and server implementations, and integration with the Energy Manager's D-Bus services.

## Features

### Core Functionality
- **REST API Server**: Unix socket-based REST API with example endpoints
- **MQTT Integration**: Subscribe and publish to Energy Manager topics
- **GDR Data Handling**: Process smart meter values using the Global Data Record (GDR) format
- **Frontend Visualization**: Display backend-provided GDR smart meter values in a tabular frontend view
- **Signal Handling**: Graceful shutdown on SIGINT/SIGTERM

### Modbus Support
- **Modbus TCP Server**: Listen on port 502 with configurable holding registers and coils
- **Modbus RTU Server**: Serial interface support via `/dev/ttyAPP1`
- **Modbus Client**: Support for both TCP and RTU protocols with read/write operations
- **MQTT-Modbus Bridge**: Control and query Modbus devices via MQTT messages

### Frontend
- **Frontend**: Demo GUI that renders GDR data from the backend as a table to illustrate frontend integration

## Prerequisites

- TQ-Systems Energy Manager
- Go 1.23.6 or later (for development)
- Access to Energy Manager's MQTT broker
- Network access for Modbus TCP (port 502 is configured in firewall)
- Yarn 4.6.0

## Building

The application uses the TQ Energy Manager Application SDK build system to compile the Go backend/frontend and package it according to the Energy Manager's app specifications. 

## Installation

Deploy the built application to the Energy Manager following the standard TQ Energy Manager app installation procedure.

## Usage

### Command Line Options

```bash
em-app-go-demo [options]
```

**Options:**
- `-loglevel <level>`: Set log level (debug, info, warning, error, panic, fatal). Default: `info`
- `-logconsole`: Write logs to STDOUT instead of syslog
- `-version`: Show version information
- `-broker <host>`: MQTT broker host. Default: `127.0.0.1`
- `-broker-port <port>`: MQTT broker port. Default: `1883`
- `-listen <path>`: REST API listening address. Default: `/run/em/apps/go-demo/socket`
- `-listenprotocol <protocol>`: REST API protocol. Default: `unix`
- `-listengroup <group>`: User group for unix socket. Default: `www`

### REST API Endpoints

The application exposes a REST API on a Unix socket (default: `/run/em/apps/go-demo/socket`).

The frontend uses the backend integration together with GDR websocket data to present current smart meter values in a table.

**Base URL:** `/api/go-demo`

#### GET /time
Returns the current server time in JSON format.

**Example:**
```bash
curl --unix-socket /run/em/apps/go-demo/socket http://localhost/api/go-demo/time
```

**Response:**
```json
"2025-11-17T10:30:45.123456789Z"
```

### MQTT Topics

#### Subscribed Topics

**Smart Meter Values** (`gdr/local/values/smart-meter`):
- Receives GDR-encoded smart meter data (voltage, current, power, etc.)
- Automatically processed and logged by the application

**Modbus Client Control** (`modbus/local/values/client/in`):
- Send Modbus read/write requests via MQTT
- Supports TCP and RTU protocols

**Modbus Server Control** (`modbus/local/values/server/in`):
- Read or write Modbus server registers via MQTT
- Update holding registers dynamically

#### Published Topics

**Modbus Client Responses** (`modbus/local/values/client/out`):
- Results from Modbus client operations

**Modbus Server Updates** (`modbus/local/values/server/out`):
- Notifications when server registers are modified

### Modbus Client Usage

Send MQTT messages to `modbus/local/values/client/in` to perform Modbus operations.

#### TCP Read Example
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 0,
  "quantity": 4,
  "function": "holding"
}
```

#### RTU Read Example
```json
{
  "protocol": "rtu",
  "device": "/dev/ttyAPP1",
  "baud": 19200,
  "dataBits": 8,
  "stopBits": 1,
  "parity": "N",
  "unitId": 1,
  "start": 0,
  "quantity": 2,
  "function": "input"
}
```

#### Write Holding Registers
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 0,
  "quantity": 1,
  "function": "holding",
  "op": "write",
  "values": [42]
}
```

#### Write Coils
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 0,
  "quantity": 1,
  "function": "coil",
  "op": "write",
  "bits": [true]
}
```

**Supported Functions:**
- `coil` - Coils (FC 1 read, FC 5/15 write)
- `discrete` - Discrete Inputs (FC 2 read)
- `holding` - Holding Registers (FC 3 read, FC 6/16 write)
- `input` - Input Registers (FC 4 read)

### Modbus Server Usage

The application runs a Modbus TCP server on `0.0.0.0:502` and optionally a Modbus RTU server on `/dev/ttyAPP1`.

#### Default Register Values

**Holding Registers:**
- Address 0: `42`
- Address 1: `0xdead`
- Address 2: `0xbeef`
- Address 3: Auto-incrementing counter (updates every second)

**Coils:**
- Address 0: `1` (ON)

#### Access via Standard Modbus Client

You can connect to the server using any Modbus TCP client:

```python
from pymodbus.client import ModbusTcpClient

client = ModbusTcpClient('energy-manager-ip', port=502)
result = client.read_holding_registers(0, 4, unit=1)
print(result.registers)  # [42, 57005, 48879, <counter>]
```

#### Access via MQTT

**Read Register:**
```json
{
  "addr": 3
}
```

**Write Register:**
```json
{
  "addr": 0,
  "value": 100
}
```

Publish to `modbus/local/values/server/in`. The server will respond on `modbus/local/values/server/out`.

## Configuration

### Firewall Rules

The application requires TCP port 502 (Modbus) to be open. This is configured in `em-fw.conf`:

```
tcp dport {502} accept
```

### Modbus RTU Settings

To modify Modbus RTU parameters, edit the constants in `backend/modbus/defines.go`:

```go
DefaultModbusRtuBaudRate = 19200
DefaultModbusRtuDataBits = 8
DefaultModbusRtuStopBits = 1
DefaultModbusRtuParity   = "N"
DefaultModbusServerRtuInterfaceName = "APP1"
EnableModbusServerRtu = true
```

## Development

### Project Structure

```
.
├── backend/
│   ├── main.go                     # Application entry point
│   ├── go.mod                      # Go module dependencies
│   └── modbus/                     # Modbus client/server implementation
│       ├── client.go               # Modbus client with MQTT control
│       ├── server.go               # Modbus TCP/RTU server
│       ├── defines.go              # Constants and data structures
│       └── util.go                 # Helper functions
├── frontend/
|   ├── src 
|       ├── components
|           ├── CardDatapoints.vue  # displays backend-provided GDR smart meter values in a table
|       ├── lang/default
|           ├── de.json             # german translation strings
|           ├── en.json             # english translation strings
|       ├── utils
|           ├── gdr.ts              # GDR WebSocket helpers and type re-exports
|           ├── obis.ts             # OBIS code parsing and decoding utilities
|       ├── ViewApp.vue             # root app-template <- start your frontend here
├── docs/                           # Documentation
├── Makefile                        # Build configuration
├── em-fw.conf                      # Firewall configuration
└── README.md                       # This file
```

### Key Dependencies

- `github.com/tq-systems/em-gdr/v2` - Global Data Record format
- `github.com/tq-systems/go-dbus` - D-Bus serial interface
- `github.com/tq-systems/public-go-utils/v3` - MQTT, REST, logging utilities
- `github.com/goburrow/modbus` - Modbus client implementation
- `github.com/tbrandon/mbserver` - Modbus server implementation

### Adding Custom Functionality

1. **Add REST Endpoints**: Define new routes in the `routes` slice in `main.go`
2. **Subscribe to MQTT Topics**: Use `mqttClient.Subscribe()` with a callback function
3. **Extend Modbus Functionality**: Modify the `modbus` package or add custom function handlers

### Example: Adding a Custom REST Endpoint

```go
routes := []rest.Route{
    {Method: "GET", Pattern: "/time", Role: "user", Handler: handleGetTimeRequest},
    {Method: "GET", Pattern: "/custom", Role: "user", Handler: handleCustomRequest},
}

func handleCustomRequest(r *http.Request) *rest.Response {
    return rest.NewJSONResponse(map[string]string{"message": "Hello!"})
}
```

## Logging

Logs are written to syslog by default. Use `-logconsole` flag to redirect to STDOUT for debugging.

**Log Levels:**
- `debug`: Detailed information for debugging
- `info`: General informational messages (default)
- `warning`: Warning messages
- `error`: Error messages
- `panic`: Panic-level messages
- `fatal`: Fatal errors that cause application exit

## Troubleshooting

### Modbus Server Not Accessible
- Verify port 502 is open: `sudo nft list ruleset | grep 502`
- Check if the application is running: `ps aux | grep go-demo`
- Review logs: `journalctl -u em-app-go-demo`

### MQTT Connection Issues
- Verify broker is running: `systemctl status mosquitto`
- Check broker host/port configuration
- Ensure MQTT client ID is unique

### Serial Interface Not Available
- Check if `/dev/ttyAPP1` exists
- Verify serial interface permissions
- Review D-Bus serial service status

## License Information

All files in this project are classified as product-specific software and bound to the use with the TQ-Systems GmbH product: TQ Energy Manager Application SDK

    SPDX-License-Identifier: LicenseRef-TQSPSLA-1.0.3

Copyright (c) 2026 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
