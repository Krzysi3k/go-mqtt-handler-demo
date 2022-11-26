package iot

import (
	"context"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v9"
)

type Event struct {
	topic      string
	msg        string
	mqttClient MQTT.Client
	ctx        context.Context
	rdb        *redis.Client
}

var (
	funcMapping = map[string]func(event Event){
		"zigbee2mqtt/SypialniaButton":      sypialniaHandleEvent,
		"zigbee2mqtt/Termometr":            termometrHandleEvent,
		"zigbee2mqtt/VibrationSensor":      storeSensorData,
		"zigbee2mqtt/MotionSensor":         storeMotionData,
		"zigbee2mqtt/SmartCube":            smartcubeHandleEvent,
		"zigbee2mqtt/SalonButton":          salonButtonHandleEvent,
		"zigbee2mqtt/DrzwiWejscioweButton": doorButtonHandleEvent,
		"zigbee2mqtt/DrzwiWejscioweSensor": doorSensorHandleEvent,
		"zigbee2mqtt/MultiButton":          multiButtonHandleEvent,
		"zigbee2mqtt/Knob":                 knobHandleEvent,
		"cmnd/zigbee2mqtt/ZarowkaRgb":      bulbSypialniaHandleEvent,
		"cmnd/DrzwiWejscioweSensor/state":  getDoorSensorState,
		"cmnd/nodemcu/Telewizor":           tvHandleEvent,
		"cmnd/pingstatus/Tv":               pingStatus,
		"cmnd/pingstatus/Ryzen":            pingStatus,
		"cmnd/pingstatus/Hypervisor":       pingStatus,
		"cmnd/rpi/docker_info":             dockerInfo,
		"cmnd/rpi/redis_info":              redisInfo,
		"cmnd/hypervisor/WakeOnLan":        wakeOnLan,
		"tele/gniazdko-4/SENSOR":           washingMachineHandleEvent,
	}
)

func (event Event) Dispatch() {
	logEvent(event.topic, event.msg)
	fn := funcMapping[event.topic]
	if fn != nil {
		fn(event)
	}
}

func NewEvent(topic, msg string, c MQTT.Client, ctx context.Context, rdb *redis.Client) Event {
	return Event{
		topic:      topic,
		msg:        msg,
		mqttClient: c,
		ctx:        ctx,
		rdb:        rdb,
	}
}
