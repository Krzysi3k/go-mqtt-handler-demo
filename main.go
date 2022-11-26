package main

import (
	"context"
	"fmt"
	"mqtt-handler-go/addr"
	"mqtt-handler-go/config"

	"github.com/go-redis/redis/v9"
	guid "github.com/google/uuid"
)

func main() {

	id := guid.New().String()
	clientID := fmt.Sprintf(addr.ClientName + "-" + id)
	ctx := context.Background()

	r := redis.NewClient(&redis.Options{
		Addr:     addr.WyseIP + ":6379",
		Password: "",
		DB:       0,
	})

	mqttCfg := config.Config{
		BrokerAddr: addr.BrokerConn,
		ClientID:   clientID,
		Topic:      "+/+/#",
		Ctx:        ctx,
		Rdb:        r,
	}

	mqttCfg.InitMqtt()
}
