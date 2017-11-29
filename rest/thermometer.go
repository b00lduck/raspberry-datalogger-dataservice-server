package rest

import (
    "github.com/gocraft/web"
)

// Get thermometer reading by code
func (c *Context) ThermometerByCodeHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleRead(rw, req)
}

// Add thermometer reading by code
func (c *Context) ThermometerByCodeAddReadingHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleWrite(rw, req)
}
