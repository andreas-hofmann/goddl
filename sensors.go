package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type SensorError struct {
	E string
}

func (e SensorError) Error() string {
	return fmt.Sprintf("Sensor error: %s", e.E)
}

type Sensor interface {
	LastUpdated() string
	Value() string
}

type SensorHumidity struct {
	Humidity    int
	Lastupdated string
}

func (s SensorHumidity) String() string {
	return fmt.Sprintf("Humidity: %d\n", s.Humidity)
}
func (s SensorHumidity) LastUpdated() string {
	return fmt.Sprintf("%s", s.Lastupdated)
}
func (s SensorHumidity) Value() string {
	return fmt.Sprintf("%.02f%%", float64(s.Humidity)/100.0)
}

type SensorTemperature struct {
	Temperature int
	Lastupdated string
}

func (s SensorTemperature) String() string {
	return fmt.Sprintf("Temperature: %d", s.Temperature)
}
func (s SensorTemperature) LastUpdated() string {
	return fmt.Sprintf("%s", s.Lastupdated)
}
func (s SensorTemperature) Value() string {
	return fmt.Sprintf("%.02f°C", float64(s.Temperature)/100.0)
}

type SensorPressure struct {
	Pressure    int
	Lastupdated string
}

func (s SensorPressure) String() string {
	return fmt.Sprintf("Pressure: %d", s.Pressure)
}
func (s SensorPressure) LastUpdated() string {
	return fmt.Sprintf("%s", s.Lastupdated)
}
func (s SensorPressure) Value() string {
	return fmt.Sprintf("%dhPA", s.Pressure)
}

type SensorInfo struct {
	Etag             string
	Manufacturername string
	Modelid          string
	Name             string
	State            json.RawMessage
	Swversion        string
	Type             string
	Uniqueid         string
}

func (si SensorInfo) String() string {
	s, err := si.Sensor()
	var value string
	if err == nil {
		value = s.Value()
	}
	return fmt.Sprintf("Name:\t%s\tType:\t%s\tValue:\t%s", si.Name, si.Type, value)
}
func (si SensorInfo) Info() string {
	return fmt.Sprintf("Name:\t%s\tType:\t%s", si.Name, si.Type)
}
func (si SensorInfo) Sensor() (Sensor, error) {
	switch si.Type {
	case "ZHATemperature":
		var temp SensorTemperature
		json.Unmarshal(si.State, &temp)
		return temp, nil
	case "ZHAHumidity":
		var hum SensorHumidity
		json.Unmarshal(si.State, &hum)
		return hum, nil
	case "ZHAPressure":
		var pres SensorPressure
		json.Unmarshal(si.State, &pres)
		return pres, nil
	default:
		return nil, SensorError{"Unknown sensor"}
	}
}
func (si SensorInfo) SensorValue() string {
	s, err := si.Sensor()
	if err != nil {
		return "Unknown"
	} else {
		return s.Value()
	}
}

type SensorMap map[int]SensorInfo

func (sm *SensorMap) HeaderStrings() []string {
	var result []string
	result = append(result, "Timestamp")
	for i := 1; i <= len(*sm); i++ {
		result = append(result, (*sm)[i].Name)
	}
	return result
}

func (sm *SensorMap) DataStrings() []string {
	var result []string
	result = append(result, time.Now().Format("2006-01-02 15:04:05"))
	for i := 1; i <= len(*sm); i++ {
		result = append(result, (*sm)[i].SensorValue())
	}
	return result
}
