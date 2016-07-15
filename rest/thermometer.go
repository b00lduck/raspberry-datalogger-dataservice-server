package rest

import (
	"github.com/gocraft/web"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/orm"
	"net/http"
	"io/ioutil"
	"strconv"
)

// Get all thermometers
func (c *Context) ThermometerHandler(rw web.ResponseWriter, req *web.Request) {
	var thermometers []orm.Counter
	db.Find(&thermometers)
	marshal(rw, thermometers)
}

// Get specific thermometer by code
func (c *Context) ThermometerByCodeHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var thermometer orm.Thermometer
	db.Where(&orm.Thermometer{Code: code}).First(&thermometer)

	if (thermometer.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Thermometer not found"))
		return
	}

	marshal(rw, thermometer)
}

// Add thermometer reading by code
func (c *Context) ThermometerByCodeAddReadingHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var thermometer orm.Thermometer
	db.Where(&orm.Thermometer{Code: code}).First(&thermometer)

	if (thermometer.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Thermometer not found"))
		return
	}

	body, err := ioutil.ReadAll(req.Body);
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not read body"))
		return
	}

	reading, err := strconv.Atoi(string(body))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Could not parse reading"))
		return
	}

	thermometerReading := orm.NewThermometerReading(thermometer, uint64(reading))
	db.Create(&thermometerReading)

	thermometer.Reading = uint64(reading)
	thermometer.LastReading = thermometerReading.Timestamp
	db.Save(thermometer)

	marshal(rw, thermometerReading)
}

// Get thermometer readings in a optionally given time range
// Query parameters: start,end
func (c *Context) ThermometerByCodeGetReadingsHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var thermometer orm.Thermometer
	db.Where(&orm.Thermometer{Code: code}).First(&thermometer)

	if (thermometer.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Thermometer not found"))
		return
	}

	start,err := c.parseUintQueryParameter(rw, "start")
	if (err != nil) {
		return
	}

	end,err := c.parseUintQueryParameter(rw, "end")
	if (err != nil) {
		return
	}

	var thermometerReadings []orm.ThermometerReading
	orm.GetOrderedWindowedQuery(db, "thermometer_id", thermometer.ID, start, end).Find(&thermometerReadings)

	startReading := orm.ThermometerReading{
		Reading: thermometerReadings[0].Reading,
		Timestamp: start,
		ThermometerID: thermometer.ID,
	}

	endReading := orm.ThermometerReading{
		Reading: thermometerReadings[len(thermometerReadings)-1].Reading,
		Timestamp: end,
		ThermometerID: thermometer.ID,
	}

	thermometerReadings = append([]orm.ThermometerReading{startReading}, thermometerReadings...)
	thermometerReadings = append(thermometerReadings, endReading)

	marshal(rw, thermometerReadings)
}