package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
        log "github.com/Sirupsen/logrus"
	"os"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/orm"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/initialization"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/rest"
)

func main() {

        log.Info("Starting dataservice...");

	mysql := os.Getenv("MYSQL_HOST")

	db, err := gorm.Open("mysql", "root:root@tcp(" + mysql + ")/rem-dataservice?parseTime=true")
	if err != nil {
	    log.WithField("error", err).Fatal("Could not initialize database")
	}

	db.SingularTable(true)

	db.AutoMigrate(&orm.Counter{}, &orm.CounterEvent{},
		&orm.Thermometer{}, &orm.ThermometerReading{},
		&orm.Flag{}, &orm.FlagState{})

	db.Model(&orm.CounterEvent{}).AddForeignKey("counter_id", "counter(id)", "RESTRICT", "RESTRICT")
	db.Model(&orm.ThermometerReading{}).AddForeignKey("thermometer_id", "thermometer(id)", "RESTRICT", "RESTRICT")
	db.Model(&orm.FlagState{}).AddForeignKey("flag_id", "flag(id)", "RESTRICT", "RESTRICT")

        log.Info("Initializing tables...")
	counterChecker := initialization.NewCounterChecker(db)
	counterChecker.CheckCounters()

        log.Info("Starting REST server...")
	rest.StartServer(db)
}


