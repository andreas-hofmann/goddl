package main

import (
	"encoding/csv"
	"os"
	"reflect"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/withmandala/go-log"
)

func main() {
	var args struct {
		Ip          string
		Apikey      string
		Logfile     string `default:"./data.log"`
		Configfile  string `default:"./config.yaml"`
		StoreConfig bool   `default:"false"`
		Debug       bool   `default:"false"`
		Logtype     string `default:"csv" help:"The data format to log. Possible values: json,csv"`
	}

	arg.MustParse(&args)

	logger := log.New(os.Stderr)

	if args.Debug {
		logger.WithDebug()
	}

	var config Config = NewConfig()

	err := config.Read(args.Configfile)
	if err != nil {
		logger.Debug(err)
	}

	if args.Apikey != "" {
		config.ApiKey = args.Apikey
	}

	if args.Ip != "" {
		config.RemoteIp = args.Ip
	}

	if config.RemoteIp == "" {
		logger.Fatal("Missing IP.")
	}

	if args.Apikey != "" {
		config.ApiKey = args.Apikey
	}

	if args.Logtype != "" {
		config.Logtype = args.Logtype
	}

	switch config.Logtype {
	case "csv":
		logger.Debug("Logging CSV data.")
	case "json":
		logger.Debug("Logging JSON data.")
	default:
		logger.Fatal("Invalid log format:", args.Logtype)
	}

	if config.ApiKey == "" {
		logger.Info("No API key supplied. Registering one.")
		answer, err := Register(config.RemoteIp)
		if err != nil {
			logger.Fatal("Error registering apikey:", err)
		} else if answer.Error.Type != 0 {
			logger.Fatal("Error registering apikey:", answer.Error.Description)
		} else {
			logger.Info("Acquired apikey:", answer.Success.Username)
			config.ApiKey = answer.Success.Username
		}
	}

	if args.StoreConfig {
		if config.Write(args.Configfile) != nil {
			logger.Warn("Error saving config to", args.Configfile)
		}
	}

	var lastSensors SensorMap

	logger.Info("Connecting to host", config.RemoteIp)

	logfile, err := os.OpenFile(config.Logfile, (os.O_WRONLY | os.O_APPEND), 0644)
	if err != nil {
		logfile, err = os.Create(config.Logfile)
		if err != nil {
			logger.Fatal("Error creating logfile")
		}
	}
	defer logfile.Close()

	var headerWritten bool = false
	csvWriter := csv.NewWriter(logfile)

	for {
		sensors, error := GetSensorMap(config.RemoteIp, config.ApiKey)
		if error != nil {
			logger.Fatal("Error fetching sensor list:", error)
		}

		if !reflect.DeepEqual(lastSensors, sensors) {
			if args.Logtype == "csv" {
				if !headerWritten {
					if csvWriter.Write(sensors.HeaderStrings()) != nil {
						logger.Fatal("Error writing header to logfile")
					}
					headerWritten = true
					logger.Debug("Header ->", sensors.HeaderStrings())
				}

				if csvWriter.Write(sensors.DataStrings()) != nil {
					logger.Fatal("Error writing data to logfile")
				}

				csvWriter.Flush()

				logger.Debug("Data   ->", sensors.DataStrings())
			} else if args.Logtype == "json" {
				data, err := sensors.Json(false)
				if err != nil {
					logger.Fatal("Error writing data to logfile")
				}

				logfile.Write(data)

				logger.Debug("Read data:", string(data))
			}
		}

		lastSensors = sensors

		time.Sleep(time.Duration(config.PollInterval) * time.Second)
	}
}
