package rest

import (
	"github.com/gocraft/web"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/orm"
	"net/http"
	"io/ioutil"
	"fmt"
)

// Get all counters
func (c *Context) CounterHandler(rw web.ResponseWriter, req *web.Request) {
	var counters []orm.Counter
	db.Find(&counters)
	marshal(rw, counters)
}

// Get specific counter by code
func (c *Context) CounterByCodeHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var counter orm.Counter
	db.Where(&orm.Counter{Code: code}).First(&counter)

	if (counter.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Counter not found"))
		return
	}

	marshal(rw, counter)
}

// Tick counter by code
func (c *Context) CounterByCodeTickHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var counter orm.Counter
	db.Where(&orm.Counter{Code: code}).First(&counter)

	if (counter.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Counter not found"))
		return
	}

	counterEvent := orm.NewTickCounterEvent(counter)
	db.Create(&counterEvent)

	counter.Reading = counterEvent.Reading
	counter.LastTick = counterEvent.Timestamp
	db.Save(counter)

	marshal(rw, counterEvent)
}

// Correct counter by code
func (c *Context) CounterByCodeCorrectHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var counter orm.Counter
	db.Where(&orm.Counter{Code: code}).First(&counter)

	if (counter.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Counter not found"))
		return
	}

	hah, err := ioutil.ReadAll(req.Body);
	fmt.Println(hah)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error reading body"))
		return
	}

	newReading,err := parseUintFromString(string(hah))
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Malformed value"))
		return
	}

	delta := int64(newReading) - int64(counter.Reading)
	counter.Reading = newReading
	db.Save(counter)

	counterEvent := orm.NewAbsCorrCounterEvent(counter, delta)

	db.Create(&counterEvent)
	marshal(rw, counterEvent)
}

// Get counter events in a optionally given time range
// Query parameters: start,end
func (c *Context) CounterByCodeEventsHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var counter orm.Counter
	db.Where(&orm.Counter{Code: code}).First(&counter)

	if (counter.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Counter not found"))
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

	var counterEvents []orm.CounterEvent
	orm.GetOrderedWindowedQuery(db, "counter_id", counter.ID, start, end).Find(&counterEvents)

	lastEvent := orm.NewLastCounterEvent(counter)
	if lastEvent.Timestamp > end {
		counterEvents = append(counterEvents, lastEvent)
	}

	marshal(rw, counterEvents)
}
