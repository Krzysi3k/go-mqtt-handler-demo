package iot

import (
	"strconv"
	"strings"
)

func bulbSypialniaHandleEvent(event Event) {
	if strings.Contains(event.msg, "#") {
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"color":{"hex":"`+event.msg+`"}}`)
	} else {
		val, err := strconv.ParseFloat(event.msg, 64)
		if err != nil {
			logError(err, "cannot parse to float64"+event.msg)
			return
		}
		multiplied := 2.54 * val
		jsonVal := strconv.FormatFloat(multiplied, 'f', 2, 64)
		event.mqttClient.Publish("zigbee2mqtt/ZarowkaRgb/set", 0, false, `{"brightness": `+jsonVal+`,"transition":5}`)

        // cache those values in redis:
        // todo 
	}
}
