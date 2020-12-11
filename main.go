package main

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		Ip          string
		Apikey      string
		Logfile     string `default:"./log.csv"`
		Configfile  string `default:"./config.yaml"`
		StoreConfig bool   `default:"false"`
	}

	arg.MustParse(&args)

	var config Config = NewConfig()

	err := config.Read(args.Configfile)
	if err != nil {
		log.Println(err)
	}

	if args.Apikey != "" {
		config.ApiKey = args.Apikey
	}

	if args.Ip != "" {
		config.RemoteIp = args.Ip
	}

	if config.RemoteIp == "" {
		log.Fatal("Missing IP.")
	}

	if args.Apikey != "" {
		config.ApiKey = args.Apikey
	}

	if config.ApiKey == "" {
		log.Print("No API key supplied. Registering one.")
		answer, err := Register(config.RemoteIp)
		if err != nil {
			log.Fatal("Error registering apikey: ", err)
		} else if answer.Error.Type != 0 {
			log.Fatal("Error registering apikey: ", answer.Error.Description)
		} else {
			log.Println("Acquired apikey:", answer.Success.Username)
			config.ApiKey = answer.Success.Username
		}
	}

	if args.StoreConfig {
		if config.Write(args.Configfile) != nil {
			log.Print("Error saving config to ", args.Configfile)
		}
	}

	var lastSensors SensorMap

	log.Println("Connecting to host ", config.RemoteIp)

	logfile, err := os.OpenFile(config.Logfile, (os.O_WRONLY | os.O_APPEND), 0644)
	if err != nil {
		logfile, err = os.Create(config.Logfile)
		if err != nil {
			log.Panic("Error creating logfile")
		}
	}
	defer logfile.Close()

	var headerWritten bool = false
	csvWriter := csv.NewWriter(logfile)

	for {
		sensors, error := GetSensorMap(config.RemoteIp, config.ApiKey)
		if error != nil {
			log.Fatal("Error fetching sensor list:", error)
		}

		if !reflect.DeepEqual(lastSensors, sensors) {
			if !headerWritten {
				if csvWriter.Write(sensors.HeaderStrings()) != nil {
					log.Fatal("Error writing header to logfile")
				}
				headerWritten = true
			}

			if csvWriter.Write(sensors.DataStrings()) != nil {
				log.Fatal("Error writing data to logfile")
			}

			csvWriter.Flush()
		}

		lastSensors = sensors

		time.Sleep(time.Duration(config.PollInterval) * time.Second)
	}
}
