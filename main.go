package main

import (
	"encoding/csv"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/withmandala/go-log"
)

const tabsize float64 = 8.0

func printSnap(logger *log.Logger, sensors SensorMap) {
	maxlen := 0
	for i := 1; i <= len(sensors); i++ {
		if sensors[i].Valid() {
			length := len(sensors[i].Name)
			if length > maxlen {
				maxlen = length
			}
		}
	}

	for i := 1; i <= len(sensors); i++ {
		if sensors[i].Valid() {
			separators := ""
			maxTabs := int(math.Ceil(float64(maxlen)) / tabsize)
			nrTabs := int(math.Ceil(float64(len(sensors[i].Name))) / tabsize)

			for i := nrTabs; i <= maxTabs; i++ {
				separators += "\t"
			}

			logger.Info("Found sensor", sensors[i].Name, separators, sensors[i].Type[3:], ":", sensors[i].SensorValue())
		}
	}
}

func main() {
	var args struct {
		Remote      string `help:"The gateway IP to connect to."`
		Apikey      string `help:"The APIkey to use for the gateway connection. If missing, a new one is registered."`
		Logfile     string `help:"The file to store the data. Default: ./data.log"`
		Configfile  string `help:"The configfile to use." default:"./config.yaml"`
		StoreConfig bool   `help:"Store the current config (+ a new API key) to the configfile. Default: false"`
		Debug       bool   `help:"Enable debug messages. Default: false"`
		Logtype     string `help:"The data format to log. Possible values: json,csv. Default: csv"`
		Snap        bool   `help:"If set, only poll the current sensor state once, and output the values on the terminal."`
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

	if args.Remote != "" {
		config.Remote = args.Remote
	}

	if args.Logfile != "" {
		config.Logfile = args.Logfile
	}

	if args.Logtype != "" {
		config.Logtype = args.Logtype
	}

	if config.Remote == "" {
		logger.Fatal("Missing IP.")
	}

	switch config.Logtype {
	case "csv":
		logger.Debug("Logging CSV data.")
	case "json":
		logger.Debug("Logging JSON data.")
	default:
		logger.Fatal("Invalid log format:", config.Logtype)
	}

	if config.ApiKey == "" {
		logger.Info("No API key supplied. Registering one.")
		answer, err := Register(config.Remote)
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

	logger.Info("Connecting to host", config.Remote)

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

Pollloop:
	for {
		sensors, error := GetSensorMap(config.Remote, config.ApiKey)
		if error != nil {
			logger.Fatal("Error fetching sensor list:", error)
		}

		if args.Snap {
			printSnap(logger, sensors)
			break Pollloop
		}

		if !reflect.DeepEqual(lastSensors, sensors) {
			switch config.Logtype {
			case "csv":
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
			case "json":
				data, err := sensors.Json(false)
				if err != nil {
					logger.Fatal("Error writing data to logfile")
				}

				logfile.Write(data)

				logger.Debug("Read data:", string(data))
			default:
				logger.Fatal("Unsupported logtype:", config.Logtype)
			}

			lastSensors = sensors
		}

		time.Sleep(time.Duration(config.PollInterval) * time.Second)
	}
}
