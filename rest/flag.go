package rest

import (
	"github.com/gocraft/web"
)


// Get flag state by code
func (c *Context) FlagByCodeHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleRead(rw, req)
}

// Change flag state by code
func (c *Context) FlagByCodeChangeStateHandler(rw web.ResponseWriter, req *web.Request) {
    c.simpleWrite(rw, req)
}