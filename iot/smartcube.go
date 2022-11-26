package iot

import (
	"encoding/json"
	"fmt"
	"mqtt-handler-go/addr"
)

type SmartCube struct {
	Action            string  `json:"action"`
	Battery           float64 `json:"battery"`
	Current           float64 `json:"current"`
	DeviceTemperature float64 `json:"device_temperature"`
	Linkquality       float64 `json:"linkquality"`
	Power             float64 `json:"power"`
	PowerOutageCount  float64 `json:"power_outage_count"`
	Temperature       float64 `json:"temperature"`
	Voltage           float64 `json:"voltage"`
}

func smartcubeHandleEvent(event Event) {
	var sc SmartCube
	if err := json.Unmarshal([]byte(event.msg), &sc); err != nil {
		logError(err, "cannot unmarshal json "+event.msg)
		return
	}
	payload, err := event.rdb.Get(event.ctx, "rotate-option").Result()
	if err != nil {
		event.rdb.Set(event.ctx, "rotate-option", "volume", 0)
		return
	}

	switch sc.Action {
	case "rotate_left":
		httpGet(fmt.Sprintf("http://%v/tv/%v/down", addr.NodeMcuIP, payload))
	case "rotate_right":
		httpGet(fmt.Sprintf("http://%v/tv/%v/up", addr.NodeMcuIP, payload))
	case "shake":
		val := "volume"
		if payload == "volume" {
			val = "prog"
		}
		event.rdb.Set(event.ctx, "rotate-option", val, 0)
	case "flip180":
		event.msg = "ON"
		wakeOnLan(event)
		event.mqttClient.Publish("cmnd/hypervisor/WakeOnLan", 0, false, "OFF")
	case "flip90":
		turnAllDevices()
		sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
	case "tap":
		httpGet(fmt.Sprintf("http://%v/cm?cmnd=Power%%20TOGGLE", addr.LampaSalonIP))
		httpGet(fmt.Sprintf("http://%v/cm?cmnd=Power%%20TOGGLE", addr.LampaSwieczkiIP))
	case "slide":
		// todo
	}
}
