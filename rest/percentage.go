package rest

import (
	"github.com/gocraft/web"
)

// Get Percentage reading by code
func (c *Context) PercentageByCodeHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleRead(rw, req)
}

// Add Percentage reading by code
func (c *Context) PercentageByCodeAddReadingHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleWrite(rw, req)
}
