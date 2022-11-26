package iot

import (
	"encoding/json"
	"fmt"
	"mqtt-handler-go/addr"
)

type DoorButton struct {
	Action      string  `json:"action"`
	Battery     float64 `json:"battery"`
	Linkquality float64 `json:"linkquality"`
	Voltage     float64 `json:"voltage"`
}

type DoorSensor struct {
	Battery     float64 `json:"battery"`
	BatteryLow  bool    `json:"battery_low"`
	Contact     *bool   `json:"contact"`
	Linkquality float64 `json:"linkquality"`
	Tamper      bool    `json:"tamper"`
	Voltage     float64 `json:"voltage"`
}

func doorButtonHandleEvent(event Event) {
	var btn DoorButton
	if err := json.Unmarshal([]byte(event.msg), &btn); err != nil {
		logError(err, "cannot unmarshall data "+event.msg)
		return
	}

	switch btn.Action {
	case "single":
		httpGet(fmt.Sprintf("http://%v/cm?cmnd=Power%%20TOGGLE", addr.LampaSalonIP))
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "TOGGLE", "transition":5, "brightness": 180}`)
	case "double":
		turnAllDevices()
		sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
	case "long":
		powerOffLamps()
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"state": "OFF"}`)
	}
}

func doorSensorHandleEvent(event Event) {
	var sensor DoorSensor
	if err := json.Unmarshal([]byte(event.msg), &sensor); err != nil {
		logError(err, "cannot unmarshall data "+event.msg)
		return
	}

	if sensor.Contact != nil {
		storedState, err := event.rdb.Get(event.ctx, "door-state").Result()
		if err != nil {
			event.rdb.Set(event.ctx, "door-state", "armed", 0)
			return
		}
		if storedState == "armed" && !*sensor.Contact {
			postTelegramMsg("Uwaga drzwi otwarte!", false)
		}
	}
}

func getDoorSensorState(event Event) {
	storedState, err := event.rdb.Get(event.ctx, "door-state").Result()
	if err != nil {
		event.rdb.Set(event.ctx, "door-state", "armed", 0)
		return
	}
	switch event.msg {
	case "check":
		if storedState == "armed" {
			event.mqttClient.Publish("tele/DrzwiWejscioweSensor/state", 0, false, "ON")
		} else {
			event.mqttClient.Publish("tele/DrzwiWejscioweSensor/state", 0, false, "OFF")
		}
	case "set":
		if storedState == "armed" {
			event.rdb.Set(event.ctx, "door-state", "unarmed", 0)
			event.mqttClient.Publish("tele/DrzwiWejscioweSensor/state", 0, false, "OFF")
			postTelegramMsg("Door alarm is OFF!", false)
		} else {
			event.rdb.Set(event.ctx, "door-state", "armed", 0)
			event.mqttClient.Publish("tele/DrzwiWejscioweSensor/state", 0, false, "ON")
			postTelegramMsg("Door alarm is ON!", false)
		}
	}
}
