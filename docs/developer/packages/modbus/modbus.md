# Modbus Package

> **Disclaimer (Demo Only)**: This code is provided **solely for demonstration and educational purposes** within the TQ Energy Manager SDK. It is **not production‑ready**. Before using it in a production system you must review, harden, test, and adapt it (error handling, security, resource management, performance, observability, configuration, protocol edge cases).

This module is part of the *Go-Demo App* and is intentionally written as a **didactic demo** for SDK customers. It favors clarity and explicitness over micro‑optimizations so you can adapt it quickly for production projects.

## Overview
This package demonstrates the use of Modbus (RTU and TCP) in the Energy Manager. It provides both a lightweight Modbus **Server** (TCP + optional RTU) and a dynamic, MQTT‑driven Modbus **Client** capable of on‑demand read and write operations over TCP or RTU. It is designed for integration into the Energy Manager runtime where Modbus interactions are externally orchestrated via MQTT topics.

## Libraries
The following libraries are used for the Modbus functionalities:
- [github.com/goburrow/modbus](https://github.com/goburrow/modbus): A Go implementation of a Modbus client.
- [github.com/tbrandon/mbserver](https://github.com/tbrandon/mbserver): A Go implementation of a Modbus server.

## Modbus Server
The implementation of the Modbus server can be found in `modbus/server.go`.

### Features
- Modbus TCP server listening on `0.0.0.0:502` (change it to your needs)
- Optional Modbus RTU server (serial) — enabled via `EnableModbusServerRtu` constant
- Server ticker increments Holding Register 3 every second to simulate changing values
- Example of a custom handler for Function Code 6 (Write Single Register) publishing updates over MQTT
- MQTT bridging for server register updates (in & out topics)
- Graceful shutdown with timeout fallback and last‑resort forced close for mbserver RTU edge cases

### Custom Function Code 6 Handler
As an example of a customized callback function for Modbus requests, `server.go:onWriteSingleRegister` was implemented for function code 6. It validates payload length, updates the Holding Register array, publishes an MQTT message, and returns the standard echo response (address + value).

## Shutdown Semantics
```go
server.Destructor() // stops ticker, unsubscribes MQTT, attempts graceful Modbus shutdown
```
The server uses a timeout wrapper (`closeServerWithTimeout`) around `mbserver.Close()` to prevent indefinite blocking (known upstream limitation when RTU listeners are active). On timeout it triggers a last‑resort `brutalCloseServer` using reflection/unsafe to release resources. **The forced close path is a last resort; prefer resolving upstream if possible.**

## Modbus Client
The implementation of the Modbus client can be found in `modbus/client.go`.

### Features
- Read support for coils, discrete inputs, holding registers, input registers (FC 1,2,3,4)
- Write support for single & multiple coils (FC5/15) and single & multiple holding registers (FC6/16)
- Per‑request timeout override (`timeoutMs`)
- RTU parameter overrides (baud, data bits, stop bits, parity) with sensible defaults

### Processing Model
- Incoming MQTT messages published to `modbus/local/values/client/in`
- MQTT callback enqueues raw payload (bounded channel, size 64)
- Worker goroutine dequeues and executes:
  1. Parse & validate JSON
  2. Create protocol-specific handler and connect
  3. Apply timeout (default 3s unless overridden by `timeoutMs`)
  4. Perform read or write
  5. Publish response JSON to `modbus/local/values/client/out`
- If queue is full, new inbound messages are dropped (drop-new policy) with a warning log.

### Timeouts
Override per request via `timeoutMs`. If omitted, default is 3000 ms. Applies to the underlying handler (`TCPClientHandler.Timeout` / `RTUClientHandler.Timeout`).

## RTU Specifics
When using Modbus RTU, the following points must be observed:

- Serial interface may be (re)bound per request; binding errors are logged.
- Default RTU parameters fall back to constants in `defines.go` when omitted.
- When using RTU, both the client and the server attempt to bind the corresponding serial interface (APP1, APP2) to the app (see `serialClient.BindInterface()`). The reservation of the serial interfaces can be viewed in the Energy Manager web frontend under *Device settings / Serial interfaces*.

## Firewall rules
The `em-fw.conf` file in the go-demo apps's root directory defines rules for *nftables*. Specifically, this file defines the ports (incoming and outgoing) for Modbus TCP. Therefore, if a port other than 502 is to be used for Modbus, `em-fw.conf` must be adjusted accordingly.

## Quick Start (Local TCP Demo)
Build the app and install it on the Energy Manager. Open the MQTT broker so that a connection from outside is possible.
Then publish a request and observe response.

```bash
# Example using mosquitto utilities
mosquitto_sub -t modbus/local/values/client/out -v &
mosquitto_pub -t modbus/local/values/client/in -m '{"protocol":"tcp","address":"127.0.0.1:502","unitId":1,"start":0,"quantity":4,"function":"holding"}'
```
You should see a JSON response with the demo register values.

## Configuration Knobs
The following constants can be adjusted in `defines.go` to change the behavior of the Modbus client/server:

| Constant | Purpose | Adjust For |
|----------|---------|-----------|
| `DefaultModbusServerTcpAddrPort` | TCP listen address | Alternate port / interface |
| `EnableModbusServerRtu` | Enable RTU listener | Disable if no serial hardware |
| `DefaultModbusServerRtuInterfaceName` | Serial interface label | Choose between `APP1` and `APP2` on Energy Manager platform |
| `DefaultModbusRtu*` | RTU defaults | Device specific timing |
| `DefaultModbusClientTimeout` | Per request default timeout | Network / device latency |

## MQTT Topics
Below is a list of MQTT topics used by the Modbus client/server. See `defines.go` to adjust the topic names.

| Component | Direction | Topic | Payload |
|-----------|-----------|-------|---------|
| Server | Inbound (set register) | `modbus/local/values/server/in` | `{ "addr": <uint16>, "value": <uint16> }` |
| Server | Outbound (notify change) | `modbus/local/values/server/out` | same as inbound |
| Client | Inbound (request) | `modbus/local/values/client/in` | Request JSON (see below) |
| Client | Outbound (response) | `modbus/local/values/client/out` | Response JSON (see below) |

QoS levels: server publish uses QoS 2; client responses use QoS 1 (see code).

## Server Usage
The server handles incoming requests and manages the Modbus protocol.

### Behavior
- Initial demo holding registers:
  - HR[0]=42, HR[1]=0xDEAD, HR[2]=0xBEEF
  - HR[3] increments every second.
- Writing a single holding register via Modbus function 6 triggers an MQTT publish on `modbus/local/values/server/out`.
- Other apps / external systems can inject updates by publishing to `modbus/local/values/server/in`.

### Example: Update a Register via MQTT
```json
{"addr": 10, "value": 1234}
```
Publish to `modbus/local/values/server/in` to set Holding Register 10.

## Client Usage
The client is entirely driven by MQTT request messages. Each message is processed asynchronously by a worker goroutine to avoid blocking the MQTT subscription callback.

### Request Schema (Read)
```json
{
  "protocol": "tcp",          // "tcp" | "rtu"
  "address": "192.168.1.10:502", // required for tcp
  "unitId": 1,
  "start": 0,
  "quantity": 4,
  "function": "holding",      // coil|discrete|holding|input
  "timeoutMs": 1500            // optional override
}
```
RTU example:
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

### Request Schema (Write Holding Registers)
Single register:
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 100,
  "function": "holding",
  "op": "write",
  "values": [4660]             // 0x1234
}
```
Multiple registers:
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 200,
  "function": "holding",
  "op": "write",
  "values": [1,2,3,4]
}
```
If `quantity` is provided for multi-write it must match `len(values)` or be omitted.

### Request Schema (Write Coils)
Single coil:
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 5,
  "function": "coil",
  "op": "write",
  "bits": [true]
}
```
Multiple coils:
```json
{
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "start": 32,
  "function": "coil",
  "op": "write",
  "bits": [true,false,true,true,false]
}
```

### Response Schema
```json
{
  "ok": true,
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "unitId": 1,
  "function": "holding",
  "start": 0,
  "quantity": 4,
  "data": [100,200,300,400],   // or bits: [...] for coil/discrete
  "durationMs": 12,
  "op": "read"
}
```
Error example:
```json
{
  "ok": false,
  "error": "quantity must be > 0",
  "protocol": "tcp",
  "address": "192.168.1.10:502",
  "start": 0,
  "quantity": 0,
  "unitId": 1
}
```
