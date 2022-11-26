package iot

import (
	"encoding/json"
)

type SypialniaButton struct {
	Action      string  `json:"action"`
	Battery     float64 `json:"battery"`
	Linkquality float64 `json:"linkquality"`
	Voltage     float64 `json:"voltage"`
}

func sypialniaHandleEvent(event Event) {
	var buttonPayload SypialniaButton
	if err := json.Unmarshal([]byte(event.msg), &buttonPayload); err != nil {
		logError(err, "cannot unmarshal json "+event.msg)
		return
	}
	switch buttonPayload.Action {
	case "double":
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"brightness": 80,"transition":5, "color_temp": "neutral"}`)
	case "single":
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "TOGGLE", "transition":5, "brightness": 180}`)
	case "long":
		powerOffLamps()
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "OFF"}`)
	}
}
