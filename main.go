package main

import (
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/influxSession"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/rest"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	log.Info("Starting dataservice...")
	influxDbHost := os.Getenv("INFLUXDB_HOST")
	influxDbDatabase := os.Getenv("INFLUXDB_DATABASE")
	influxDbUsername := os.Getenv("INFLUXDB_USERNAME")
	influxDbPassword := os.Getenv("INFLUXDB_PASSWORD")

	influxSessionFactory := func() influxSession.InfluxSession {
		return influxSession.NewInfluxSession(influxDbHost, influxDbDatabase, influxDbUsername, influxDbPassword)
	}

	log.Info("Starting REST server...")
	rest.StartServer(influxSessionFactory)
}
