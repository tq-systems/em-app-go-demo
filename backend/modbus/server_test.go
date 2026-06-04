/*
 * Copyright (c) 2025-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229
 * Seefeld, Germany. All rights reserved.
 * Author: Frank Hammon and the Energy Manager development team
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
	"errors"
	"testing"

	serialmock "github.com/tq-systems/em-app-go-demo/backend/modbus/mocks/serial"
	serialClient "github.com/tq-systems/go-dbus/serial"
	"go.uber.org/mock/gomock"
)

// TestNewModbusServer_ErrorListSerialInterfaces verifies correct handling
// in case the ListInterfaces function returns an error
func TestNewModbusServer_ErrorListSerialInterfaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSerial := serialmock.NewMockClient(ctrl)
	mockSerial.EXPECT().ListInterfaces().Return(nil, errors.New("TestError"))

	sut := &ModbusServer{
		serialClient: mockSerial,
	}
	err := sut.bindModbusRtu()
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if err.Error() != "unable to get device list: TestError" {
		t.Fatalf("got wrong error: %v", err)
	}
}

// TestNewModbusServer_ErrorSerialInterfaceNotListed verifies correct handling
// in case the desired interface is not listed
func TestNewModbusServer_ErrorSerialInterfaceNotListed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSerial := serialmock.NewMockClient(ctrl)
	interfaceMap := make(map[string]serialClient.BindState)
	mockSerial.EXPECT().ListInterfaces().Return(interfaceMap, nil)

	sut := &ModbusServer{
		serialClient: mockSerial,
	}
	err := sut.bindModbusRtu()
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if err.Error() != "serial interface not available" {
		t.Fatalf("got wrong error: %v", err)
	}
}

// TestNewModbusServer_ErrorSerialInterfaceUnavailable verifies correct handling
// in case the desired interface is already used
func TestNewModbusServer_ErrorSerialInterfaceUnavailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSerial := serialmock.NewMockClient(ctrl)
	interfaceMap := make(map[string]serialClient.BindState)
	interfaceMap[DefaultModbusServerRtuInterfaceName] = serialClient.BindUnavailable
	mockSerial.EXPECT().ListInterfaces().Return(interfaceMap, nil)

	sut := &ModbusServer{
		serialClient: mockSerial,
	}
	err := sut.bindModbusRtu()
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if err.Error() != "serial interface is already used by another service" {
		t.Fatalf("got wrong error: %v", err)
	}
}

// TestNewModbusServer_SerialInterfaceAlreadyBound verifies correct handling
// in case the desired interface is already bound
func TestNewModbusServer_SerialInterfaceAlreadyBound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSerial := serialmock.NewMockClient(ctrl)
	interfaceMap := make(map[string]serialClient.BindState)
	interfaceMap[DefaultModbusServerRtuInterfaceName] = serialClient.BindBound
	mockSerial.EXPECT().ListInterfaces().Return(interfaceMap, nil)

	sut := &ModbusServer{
		serialClient: mockSerial,
	}
	err := sut.bindModbusRtu()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNewModbusServer_SerialInterfaceUnbound verifies correct handling
// in case the desired interface is unbound
func TestNewModbusServer_SerialInterfaceUnbound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSerial := serialmock.NewMockClient(ctrl)
	interfaceMap := make(map[string]serialClient.BindState)
	interfaceMap[DefaultModbusServerRtuInterfaceName] = serialClient.BindUnbound
	mockSerial.EXPECT().ListInterfaces().Return(interfaceMap, nil)
	mockSerial.EXPECT().BindInterface(DefaultModbusServerRtuInterfaceName, nil, false)

	sut := &ModbusServer{
		serialClient: mockSerial,
	}
	err := sut.bindModbusRtu()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
