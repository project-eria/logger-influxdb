package main

import (
	"fmt"
	"os"

	"github.com/project-eria/eria-base"
	configmanager "github.com/project-eria/eria-base/config-manager"
	"github.com/project-eria/xaal-go"
	"github.com/project-eria/xaal-go/device"
	"github.com/project-eria/xaal-go/message"

	influxdb "github.com/influxdata/influxdb1-client/v2"
	logger "github.com/project-eria/eria-logger"
)

var (
	// Version is a placeholder that will receive the git tag version during build time
	Version = "-"
)

const configFile = "logger-influxdb.json"

func setupDev(dev *device.Device) {
	dev.VendorID = "ERIA"
	dev.ProductID = "InfluxDB Logger"
	dev.Info = fmt.Sprintf("%s@%s:%d", config.Database, config.Host, config.Port)
	dev.URL = "https://www.influxdata.com"
	dev.Version = Version
}

var config = struct {
	Addr     string
	Database string `default:"eria"`
	Host     string `default:"127.0.0.1"`
	Port     int    `default:"8086"`
	Username string `required:"true"`
	Password string `required:"true"`
	Devices  []string
}{}

var _client influxdb.Client

func main() {
	defer os.Exit(0)

	eria.AddShowVersion(Version)

	logger.Module("main").Infof("Starting InfluxDB Logger %s...", Version)

	// Loading config
	cm, err := configmanager.Init(configFile, &config)
	if err != nil {
		if configmanager.IsFileMissing(err) {
			err = cm.Save()
			if err != nil {
				logger.Module("main").WithField("filename", configFile).Fatal(err)
			}
			logger.Module("main").Fatal("JSON Config file do not exists, created...")
		} else {
			logger.Module("main").WithField("filename", configFile).Fatal(err)
		}
	}

	if err := cm.Load(); err != nil {
		logger.Module("main").Fatal(err)
	}
	defer cm.Close()

	// Init xAAL engine
	eria.InitEngine()

	setup()
	// Save for new Address during setup
	cm.Save()

	xaal.AddRxHandler(parse)

	// Launch the xAAL engine
	go xaal.Run()
	defer xaal.Stop()

	eria.WaitForExit()
}

func setup() {
	dev, err := device.New("logger.basic", config.Addr)
	if err != nil {
		logger.Module("main").WithError(err).Fatal("Error when creating the device")
	}
	if config.Addr == "" {
		config.Addr = dev.Address
	}
	setupDev(dev)
	xaal.AddDevice(dev)
	// InfluxDB
	// Create a new HTTPClient
	_client, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		logger.Module("main").Fatal(err)
	}
}

func parse(msg *message.Message) {
	if msg.IsAttributesChange() && stringInSlice(msg.Header.Source, config.Devices) {
		// Create a new points batch
		bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database:  config.Database,
			Precision: "s",
		})
		if err != nil {
			logger.Module("main").WithField("msg", msg).Error(err)
		}

		// Create a point and add to batch
		tags := map[string]string{}

		pt, err := influxdb.NewPoint(msg.Header.Source, tags, msg.Body, msg.Time())
		if err != nil {
			logger.Module("main").Fatal(err)
		}
		bp.AddPoint(pt)
		// Write the batch
		if err := _client.Write(bp); err != nil {
			logger.Module("main").WithField("msg", msg).Error(err)
		}
		logger.Module("main").WithField("msg", msg).Trace()
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
