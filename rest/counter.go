package rest

import (
    "github.com/gocraft/web"
    "fmt"
    log "github.com/sirupsen/logrus"
)

// Get counter by code
func (c *Context) CounterByCodeHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleRead(rw, req)
}

// Tick counter by code
func (c *Context) CounterByCodeTickHandler(rw web.ResponseWriter, req *web.Request) {

    code := parseStringPathParameter(req, "code")
    newVal, err := c.influxSession.IncrementSeries(code, 0.1)
    if err != nil {
        log.Error(err)
        rw.WriteHeader(500)
    } else {
        rw.WriteHeader(200)
        rw.Write([]byte(fmt.Sprintf("%0.1f", newVal)))
    }

}

// Correct counter by code
func (c *Context) CounterByCodeCorrectHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleWrite(rw, req)
}