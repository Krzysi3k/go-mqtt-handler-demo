package iot

import (
	"encoding/json"
	"mqtt-handler-go/addr"
)

type MultiButton struct {
	Action      string  `json:"action"`
	Battery     float64 `json:"battery"`
	Linkquality float64 `json:"linkquality"`
	Voltage     float64 `json:"voltage"`
}

func multiButtonHandleEvent(event Event) {
	var btn MultiButton
	if err := json.Unmarshal([]byte(event.msg), &btn); err != nil {
		logError(err, "cannot unmarshall data "+event.msg)
		return
	}
	switch btn.Action {
	case "1_single":
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "TOGGLE", "transition":5, "brightness": 180}`)
	case "1_double":
		turnAllDevices()
		sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
	case "1_hold":
		powerOffLamps()
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "OFF"}`)
	case "2_single":
		// not used
	case "2_double":
		// not used
	case "2_hold":
		// not used
	case "3_single":
		// not used
	case "3_double":
		// not used
	case "3_hold":
		// not used
	case "4_single":
		// not used
	case "4_double":
		// not used
	case "4_hold":
		// not used
	}
}
