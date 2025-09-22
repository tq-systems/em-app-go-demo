/*
 * This file is part of the go-demo application.
 * More license information can be found in the root folder.
 *
 * SPDX-License-Identifier: LicenseRef-TQSPSLA-1.0.3
 *
 * Copyright (c) 2025 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Christoph Krutz and the Energy Manager development team
 */

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gdr "github.com/tq-systems/em-gdr/v2"
	"github.com/tq-systems/public-go-utils/v3/log"
	mqttHandler "github.com/tq-systems/public-go-utils/v3/mqtt"
	"github.com/tq-systems/public-go-utils/v3/rest"
)

const (
	baseURL = "/api/go-demo"
	// smart-meter values topic ​​(voltage, current, power, etc.)
	topic = "gdr/local/values/smart-meter"
	// MqttID is the mqtt message ID of this service
	MqttID = "em-app-go-demo"
)

var (
	// Version (see https://semver.org) is set via build scripts
	Version        = ""
	logLevel       = flag.String("loglevel", "info", "Set log level (debug, info, warning, error, panic, fatal)")
	logToConsole   = flag.Bool("logconsole", false, "Write logs to STDOUT instead of syslog")
	versionFlag    = flag.Bool("version", false, "Show version")
	mqttBrokerHost = flag.String("broker", "127.0.0.1", "MQTT broker host")
	mqttBrokerPort = flag.Int("broker-port", 1883, "MQTT broker port")
	listen         = flag.String("listen", "/run/em/apps/go-demo/socket", "Set the REST API listening address")
	listenprotocol = flag.String("listenprotocol", "unix", "Set the REST API listening protocol")
	listengroup    = flag.String("listengroup", "www", "User group of unix socket")
)

func main() {
	err := run()
	if err != nil {
		log.Errorf("Failed to init application: %v", err)
		os.Exit(1)
	}
}

// run ensures that deferred calls are also executed in case of an error
func run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	flag.Parse()

	if *versionFlag {
		fmt.Println(Version)
		return nil
	}

	// setup logging
	log.InitLogger(*logLevel, *logToConsole)

	//  Do whatever you want here:
	// ....

	// setup REST handler
	routes := []rest.Route{
		{Method: "GET", Pattern: "/time", Role: "user", Handler: handleGetTimeRequest},
	}

	listener := rest.Listener{Address: *listen, Proto: *listenprotocol, Group: *listengroup}

	// setup REST server
	restServer, err := rest.NewServer(baseURL, listener, routes)
	if err != nil {
		return fmt.Errorf("could not create server: %s", err)
	}
	restErrChan := restServer.AsyncServe()

	// setup MQTT Client
	mqttClient, err := mqttHandler.NewClient(*mqttBrokerHost, *mqttBrokerPort, MqttID)
	if err != nil {
		return fmt.Errorf("failed to initialize mqtt client: %s", err)
	}
	defer mqttClient.Close()

	subscription, err := mqttClient.Subscribe(topic, onMessage)
	if err != nil {
		log.Errorf("cannot subcribe on mqtt topic %v because: %v", topic, err)
	}
	defer subscription.Unsubscribe()

	// run event loop - quit on configured signals or errors
	log.Info("SIGINT or SIGTERM to terminate the program")
	select {
	case <-sigs:
	case err := <-restErrChan:
		return fmt.Errorf("error while running REST server: %v", err)
	}

	log.Debug("Exiting...")
	return nil
}

func handleGetTimeRequest(r *http.Request) *rest.Response {
	return rest.NewJSONResponse(time.Now())
}

// example MQTT callback function
func onMessage(topic string, msg []byte) {
	var gdrData gdr.GDRs

	err := gdrData.UnmarshalVT(msg)
	if err != nil {
		log.Error("cannot unmarshal message: ", err)
	}

	// unmarshalled GDR data can be stored internally, or passed to other components as a copy via gdrData.CloneVT()
}
