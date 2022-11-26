package iot

import (
	"encoding/json"
	"fmt"
	"mqtt-handler-go/addr"
)

type SalonButton struct {
	Action      string  `json:"action"`
	Battery     float64 `json:"battery"`
	Linkquality float64 `json:"linkquality"`
	Voltage     float64 `json:"voltage"`
}

func salonButtonHandleEvent(event Event) {
	var btn SalonButton
	if err := json.Unmarshal([]byte(event.msg), &btn); err != nil {
		logError(err, "cannot unmarshal json "+event.msg)
		return
	}

	switch btn.Action {
	case "single":
		httpGet(fmt.Sprintf("http://%v/cm?cmnd=Power%%20TOGGLE", addr.LampaSalonIP))
	case "double":
		turnAllDevices()
		sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
	case "long":
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "TOGGLE", "transition":5, "brightness": 180}`)
	}
}
