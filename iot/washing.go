package iot

import (
	"encoding/json"
)

type Energy struct {
	Time  string `json:"Time"`
	Stats Stats  `json:"ENERGY"`
}

type Stats struct {
	TotalStartTime string  `json:"TotalStartTime"`
	Total          float64 `json:"Total"`
	Yesterday      float64 `json:"Yesterday"`
	Today          float64 `json:"Today"`
	Period         float64 `json:"Period"`
	Power          float64 `json:"Power"`
	ApparentPower  float64 `json:"ApparentPower"`
	ReactivePower  float64 `json:"ReactivePower"`
	Factor         float64 `json:"Factor"`
	Voltage        float64 `json:"Voltage"`
	Current        float64 `json:"Current"`
}

type StoredValues struct {
	Counter      int64   `json:"washing_counter"`
	LastValue    float64 `json:"last_value"`
	WashingState string  `json:"washing_state"`
}

func washingMachineHandleEvent(event Event) {
	var e Energy
	var s StoredValues
	if err := json.Unmarshal([]byte(event.msg), &e); err != nil {
		logError(err, "cannot unmarshall data"+event.msg)
		return
	}
	prevValue, err := event.rdb.Get(event.ctx, "washing-state").Result()
	if err != nil {
		s.WashingState = "stopped"
		payload, _ := json.Marshal(s)
		event.rdb.Set(event.ctx, "washing-state", payload, 0)
		return
	}

	if err = json.Unmarshal([]byte(prevValue), &s); err != nil {
		logError(err, "cannot unmarshall data"+event.msg)
		return
	}
	if s.WashingState == "stopped" {
		if e.Stats.Today > s.LastValue {
			s.Counter += 1
			if s.Counter >= 5 {
				logEvent(event.topic, "register: washing started")
				event.mqttClient.Publish("tele/washingmachine/state", 0, true, "ON")
				s.WashingState = "started"
				s.Counter = 0
			}
		} else {
			s.Counter = 0
		}
	} else if s.WashingState == "started" {
		if e.Stats.Today == s.LastValue {
			s.Counter += 1
			if s.Counter >= 5 {
				logEvent(event.topic, "register: washing is finished")
				event.mqttClient.Publish("tele/washingmachine/state", 0, true, "OFF")
				s.WashingState = "stopped"
				s.Counter = 0
				if err := event.rdb.Get(event.ctx, "vibration-sensor").Err(); err != nil {
					postTelegramMsg("wyjmij pranie", false)
				}
			}
		} else {
			s.Counter = 0
		}
	}

	s.LastValue = e.Stats.Today
	jsonPayload, _ := json.Marshal(s)
	event.rdb.Set(event.ctx, "washing-state", jsonPayload, 0)
}
