package main

import (
    log "github.com/Sirupsen/logrus"
	"os"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/rest"
    "github.com/b00lduck/raspberry-datalogger-dataservice-server/influxSession"
)

func main() {

    log.Info("Starting dataservice...");
    influxDbHost := os.Getenv("INFLUXDB_HOST")
    influxDbDatabase := os.Getenv("INFLUXDB_DATABASE")
    influxDbUsername := os.Getenv("INFLUXDB_USERNAME")
    influxDbPassword := os.Getenv("INFLUXDB_PASSWORD")

    influxSessionFactory := func () influxSession.InfluxSession {
        return influxSession.NewInfluxSession(influxDbHost, influxDbDatabase, influxDbUsername, influxDbPassword)
    }

    log.Info("Starting REST server...")
	rest.StartServer(influxSessionFactory)
}


