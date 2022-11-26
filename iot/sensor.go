package iot

import (
	"encoding/json"
	"time"
)

type VibrationSensor struct {
	Action            string  `json:"action"`
	Angle             float64 `json:"angle"`
	AngleX            float64 `json:"angle_x"`
	AngleXAbsolute    float64 `json:"angle_x_absolute"`
	AngleY            float64 `json:"angle_y"`
	AngleYAbsolute    float64 `json:"angle_y_absolute"`
	AngleZ            float64 `json:"angle_z"`
	Battery           float64 `json:"battery"`
	DeviceTemperature float64 `json:"device_temperature"`
	Linkquality       float64 `json:"linkquality"`
	PowerOutageCount  float64 `json:"power_outage_count"`
	Strength          float64 `json:"strength"`
	Temperature       float64 `json:"temperature"`
	Vibration         bool    `json:"vibration"`
	Voltage           float64 `json:"voltage"`
}

type MotionSensor struct {
	Battery           float64 `json:"battery"`
	DeviceTemperature float64 `json:"device_temperature"`
	Illuminance       float64 `json:"illuminance"`
	IlluminanceLux    float64 `json:"illuminance_lux"`
	Linkquality       float64 `json:"linkquality"`
	Occupancy         bool    `json:"occupancy"`
	PowerOutageCount  float64 `json:"power_outage_count"`
	Temperature       float64 `json:"temperature"`
	Voltage           float64 `json:"voltage"`
}

func storeSensorData(event Event) {
	var data VibrationSensor
	if err := json.Unmarshal([]byte(event.msg), &data); err != nil {
		logError(err, "cannot unmarshall data "+event.msg)
		return
	}
	if data.Vibration {
		event.rdb.Set(event.ctx, "vibration-sensor", event.msg, 360*time.Second)
	}
}

func storeMotionData(event Event) {
	//todo
}
