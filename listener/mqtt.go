// Copyright © 2017 thingful

package listener

import (
	"errors"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	hub "github.com/thingful/device-hub"
)

func defaultMQTTOptions(brokerAddress, clientID string) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	opts.SetClientID(clientID)

	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	return opts
}

func defaultMQTTClient(options *mqtt.ClientOptions) mqtt.Client {
	return mqtt.NewClient(options)
}

func newMQTTListener(client mqtt.Client) (*mqttlistener, error) {

	if !client.IsConnected() {
		return nil, errors.New("mqtt client is not connected")
	}

	return &mqttlistener{
		client: client,
	}, nil
}

type mqttlistener struct {
	client mqtt.Client
}

func (m *mqttlistener) NewChannel(topic string) (hub.Channel, error) {

	if topic == "" {
		return nil, errors.New("mqtt topic is empty string")
	}

	if strings.Contains(topic, "#") {
		return nil, errors.New("mqtt wildcard (#) is not allowed")
	}

	errors := make(chan error)
	out := make(chan hub.Message)

	handler := func(client mqtt.Client, msg mqtt.Message) {
		input := newHubMessage(msg.Payload(), "MQTT", msg.Topic())
		out <- input
	}

	if token := m.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return defaultChannel{out: out, errors: errors, close: func() error {
		token := m.client.Unsubscribe(topic)
		token.Wait()
		return token.Error()
	}}, nil
}

func (m *mqttlistener) Close() error {
	// TODO : set a sensible timeout
	m.client.Disconnect(1000)
	return nil
}
