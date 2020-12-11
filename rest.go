package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func Register(remote string) (RegisterResult, error) {
	client := resty.New()
	remote = fmt.Sprintf("http://%s/api", remote)
	resp, err := client.R().SetBody(`{"devicetype":"go-dsl"}`).Post(remote)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	var answer []RegisterResult
	err = json.Unmarshal(resp.Body(), &answer)
	return answer[0], err
}

func GetSensorMap(remote, apikey string) (SensorMap, error) {
	client := resty.New()
	remote = fmt.Sprintf("http://%s/api/%s/sensors", remote, apikey)
	resp, err := client.R().Get(remote)

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	var answer SensorMap
	err = json.Unmarshal(resp.Body(), &answer)
	return answer, err
}
