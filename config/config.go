package config

import (
	"context"
	"fmt"
	"log"
	"mqtt-handler-go/iot"
	"os"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v9"
)

type Config struct {
	BrokerAddr string
	ClientID   string
	Topic      string
	Ctx        context.Context
	Rdb        *redis.Client
}

func (cfg Config) InitMqtt() {
	var wg sync.WaitGroup
	wg.Add(1)

	var onConnectHandler MQTT.OnConnectHandler = func(c MQTT.Client) {
		if token := c.Subscribe(cfg.Topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		log.Printf("Connected! Subscribing topic: %v\n", cfg.Topic)
	}

	var onMessageHandler MQTT.MessageHandler = func(c MQTT.Client, mqttMessage MQTT.Message) {
		mqttTopic := mqttMessage.Topic()
		mqttPayload := string(mqttMessage.Payload())
		event := iot.NewEvent(mqttTopic, mqttPayload, c, cfg.Ctx, cfg.Rdb)
		// for debugging:
		// goNum := runtime.NumGoroutine()
		// log.Printf("goroutines COUNT: %v\n", goNum)
		go event.Dispatch()
	}

	log.SetFlags(5)

	opts := MQTT.NewClientOptions().AddBroker(cfg.BrokerAddr)
	opts.SetClientID(cfg.ClientID)
	opts.OnConnect = onConnectHandler
	opts.SetDefaultPublishHandler(onMessageHandler)

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	wg.Wait()
}
