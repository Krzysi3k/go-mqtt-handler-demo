package iot

import (
	"encoding/json"
	"fmt"
	"mqtt-handler-go/addr"
	"strings"
	"time"
)

type Telewizor struct {
	Payload string `json:"payload"`
}

func tvHandleEvent(event Event) {
	var tv Telewizor
	if err := json.Unmarshal([]byte(event.msg), &tv); err != nil {
		logError(err, "cannot unmarshal json "+event.msg)
		return
	}

	payloadPart := strings.Split(tv.Payload, ";")[0]
	switch payloadPart {
	case "all":
		turnAllDevices()
		sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
	case "channel":
		opts := strings.Split(tv.Payload, ";")[1:]
		for _, i := range opts {
			httpGet(fmt.Sprintf("http://%v/tv/%v", addr.NodeMcuIP, i))
			time.Sleep(20 * time.Millisecond)
		}
	case "volume", "prog":
		opt := strings.Split(tv.Payload, ";")[1]
		httpGet(fmt.Sprintf("http://%v/tv/%v/%v", addr.NodeMcuIP, payloadPart, opt))
	default:
		httpGet(fmt.Sprintf("http://%v/%v", addr.NodeMcuIP, tv.Payload))
	}
}

// func isTvOnline(event Event) {
// 	sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
// }
