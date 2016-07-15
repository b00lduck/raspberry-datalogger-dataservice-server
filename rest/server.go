package rest
import (
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/gocraft/web"
	"encoding/json"
	"log"
	"strconv"
	"net/url"
)

var db *gorm.DB

type Context struct {
	values *url.Values
}

func StartServer(database *gorm.DB) {

	db = database

	router := web.New(Context{}).
	Middleware(web.LoggerMiddleware).
	Middleware((*Context).QueryVarsMiddleware).
	Middleware(CorsMiddleware).

	Get ("/counter", 					(*Context).CounterHandler).
	Get ("/counter/:code", 				(*Context).CounterByCodeHandler).
	Post("/counter/:code/tick", 		(*Context).CounterByCodeTickHandler).
	Put ("/counter/:code/corr", 		(*Context).CounterByCodeCorrectHandler).
	Get ("/counter/:code/events",		(*Context).CounterByCodeEventsHandler).

	Get ("/thermometer", 				(*Context).ThermometerHandler).
	Get ("/thermometer/:code", 			(*Context).ThermometerByCodeHandler).
	Post("/thermometer/:code/reading", 	(*Context).ThermometerByCodeAddReadingHandler).
	Get ("/thermometer/:code/readings",	(*Context).ThermometerByCodeGetReadingsHandler).

	Get ("/flag",						(*Context).FlagHandler).
	Get ("/flag/:code",					(*Context).FlagByCodeHandler).
	Post("/flag/:code/state", 			(*Context).FlagByCodeChangeStateHandler).
	Get ("/flag/:code/states",			(*Context).FlagByCodeGetStatesHandler)

	e := http.ListenAndServe(":8080", router)
	if e != nil {
	    panic(e)
	}
}

func (c *Context) QueryVarsMiddleware(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	values, err := parseQueryParams(rw, r)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Malformed URL"))
		return
	}

	c.values = &values
	next(rw, r)
}

func CorsMiddleware(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	next(rw, r)
}

func marshal(rw web.ResponseWriter, v interface{}) {

	buffer, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshalling JSON")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err = rw.Write(buffer); err != nil {
		log.Println("Error writing JSON to buffer")
		rw.WriteHeader(http.StatusInternalServerError)
	}

}

func parseUintFromString(s string) (ret uint64, err error) {
	ret, err = strconv.ParseUint(s, 10, 64)
	return
}

func (c *Context) parseUintQueryParameter(rw web.ResponseWriter, name string) (ret uint64, err error) {

	s := c.values.Get(name)
	if s == "" {
		return 0, nil
	}

	ret,err = parseUintFromString(s)
	if err != nil {
		log.Println("Error parsing uint64 " + name +": " + s)
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Malformed parameter " + name))
	}

	return
}

func parseStringPathParameter(req *web.Request, name string) string {
	return req.PathParams[name]
}

func parseQueryParams(rw web.ResponseWriter, req *web.Request) (values url.Values, err error) {

	u,err := url.Parse(req.RequestURI)
	if err != nil {
		return
	}

	values,err = url.ParseQuery(u.RawQuery)

	return
}
