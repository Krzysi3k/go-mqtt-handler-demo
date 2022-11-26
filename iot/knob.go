package iot

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Knob struct {
	Action               string  `json:"action"`
	ActionStepSize       float64 `json:"action_step_size"`
	ActionTransitionTime float64 `json:"action_transition_time"`
	Battery              float64 `json:"battery"`
	Linkquality          float64 `json:"linkquality"`
	OperationMode        string  `json:"operation_mode"`
	Voltage              float64 `json:"voltage"`
}

type KnobOption struct {
	Option string  `json:"option"`
	Value  float64 `json:"value"`
}

func knobHandleEvent(event Event) {
	var knb Knob
	if err := json.Unmarshal([]byte(event.msg), &knb); err != nil {
		logError(err, "cannot unmarshall data"+event.msg)
		return
	}

	if knb.Action == "brightness_step_down" || knb.Action == "brightness_step_up" {
		opt, err := event.rdb.Get(event.ctx, "KnobOption").Result()
		if err != nil {
			event.rdb.Set(event.ctx, "KnobOption", `{"option":"brightness", "value":0}`, 0)
			return
		}
		var knbOpt KnobOption
		if err := json.Unmarshal([]byte(opt), &knbOpt); err != nil {
			logError(err, "cannot unmarshall data"+event.msg)
			return
		}
		if knbOpt.Option == "brightness" {
			if strings.Contains(knb.Action, "up") {
				knbOpt.Value = knbOpt.Value + (knb.ActionStepSize / 2.5)
				if knbOpt.Value > 254 {
					knbOpt.Value = 254
				}
			} else if strings.Contains(knb.Action, "down") {
				knbOpt.Value = knbOpt.Value - (knb.ActionStepSize / 2.5)
				if knbOpt.Value < 0 {
					knbOpt.Value = 0
				}
			}
			if knbOpt.Value > 0 {
				knbOpt.Value = math.Round(knbOpt.Value*100) / 100
			}
			event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, fmt.Sprintf("{\"brightness\":%v, \"transition\":5}", knbOpt.Value))
			currentOpt, err := json.Marshal(knbOpt)
			if err != nil {
				logError(err, "cannot marshall data"+event.msg)
				return
			}
			event.rdb.Set(event.ctx, "KnobOption", currentOpt, 0)
		}
	}
}
