package iot

import (
	"encoding/json"
	"fmt"
	"time"
)

type Termometr struct {
	Battery     float64  `json:"battery"`
	Humidity    *float64 `json:"humidity"`
	Linkquality float64  `json:"linkquality"`
	Temperature *float64 `json:"temperature"`
	Voltage     float64  `json:"voltage"`
}

func termometrHandleEvent(event Event) {
	// log.Printf("topic: %v, payload: %v\n", event.topic, event.msg)
	var t Termometr
	if err := json.Unmarshal([]byte(event.msg), &t); err != nil {
		logError(err, "cannot unmarshal json "+event.msg)
		return
	}
	if t.Humidity != nil && t.Temperature != nil {
		ts := time.Now().Unix()
		newData := fmt.Sprintf("%v,%v,%v\n", *t.Temperature, *t.Humidity, ts)
		storedTemps, err := event.rdb.Get(event.ctx, "termometr-payload").Result()
		if err != nil {
			newData = fmt.Sprintf("temp,hum,ts\n%v,%v,%v\n", *t.Temperature, *t.Humidity, ts)
		}
		currentPayload := storedTemps + newData
		err = event.rdb.Set(event.ctx, "termometr-payload", currentPayload, 0).Err()
		logError(err, "cannot set 'termometr-payload' key"+event.msg)

		threshold := 70.0
		if *t.Humidity < threshold {
			if err := event.rdb.Get(event.ctx, "humidity").Err(); err != nil {
				event.rdb.Set(event.ctx, "humidity", *t.Humidity, time.Hour*24)
				message := fmt.Sprintf("humidity: %v%% \ntemperature: %v℃", *t.Humidity, *t.Temperature)
				time.Sleep(time.Second * 10) // set timeout
				logEvent(event.topic, "sending Telegram message")
				postTelegramMsg(message, true)
				annotationText := fmt.Sprintf("humidity below %v%%, current value %v", threshold, *t.Humidity)
				logEvent(event.topic, "sending Grafana annotations")
				postGrafanaReleaseAnnotation(GrafanaAnnotation{dashboardId: 3, panelId: 14, msg: annotationText})
				// postGrafanaReleaseAnnotation(1, 4, annotationText)
			}
		}
	}
}
