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
	"errors"
	"testing"

	mqttmock "github.com/tq-systems/em-app-go-demo/backend/modbus/mocks/mqtt"
	serialmock "github.com/tq-systems/em-app-go-demo/backend/modbus/mocks/serial"
	"go.uber.org/mock/gomock"
)

// TestNewModbusClient_SubscribeAndDestructor verifies successful client creation,
// MQTT subscription setup, and cleanup via Destructor.
func TestNewModbusClient_SubscribeAndDestructor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMQTT := mqttmock.NewMockClient(ctrl)
	mockSub := mqttmock.NewMockSubscription(ctrl)
	mockSerial := serialmock.NewMockClient(ctrl)

	mockMQTT.EXPECT().
		Subscribe(MqttModbusClientInTopic, gomock.Any()).
		Return(mockSub, nil)

	c, err := NewModbusClient(mockMQTT, mockSerial)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	mockSub.EXPECT().Unsubscribe()
	c.Destructor()
}

// TestNewModbusClient_SubscribeError verifies constructor error handling when
// the MQTT subscription fails.
func TestNewModbusClient_SubscribeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMQTT := mqttmock.NewMockClient(ctrl)
	mockSerial := serialmock.NewMockClient(ctrl)
	wantErr := errors.New("subscribe failed")

	mockMQTT.EXPECT().
		Subscribe(MqttModbusClientInTopic, gomock.Any()).
		Return(nil, wantErr)

	c, err := NewModbusClient(mockMQTT, mockSerial)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if c != nil {
		t.Fatal("expected nil client on subscribe error")
	}
}
