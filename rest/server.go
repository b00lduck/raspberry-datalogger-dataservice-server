package rest
import (
	"net/http"
	"github.com/gocraft/web"
    log "github.com/Sirupsen/logrus"
	"strconv"
	"net/url"
    "github.com/b00lduck/raspberry-datalogger-dataservice-server/influxSession"
    "fmt"
    "io/ioutil"
)

var isf func() influxSession.InfluxSession

type Context struct {
    influxSession influxSession.InfluxSession
    values *url.Values
}

func StartServer(influxSessionFactory func() influxSession.InfluxSession) {

    isf = influxSessionFactory

    initSession := isf()
    err := initSession.InitDb()
    if err != nil {
        log.Fatal(err)
    }
    initSession.Close()

	router := web.New(Context{}).
	Middleware(web.LoggerMiddleware).
    Middleware((*Context).InfluxSessionMiddleware).
	Middleware((*Context).QueryVarsMiddleware).

    Get("/counter/:code", 		        (*Context).CounterByCodeHandler).
	Post("/counter/:code/tick", 		(*Context).CounterByCodeTickHandler).
	Put ("/counter/:code",      		(*Context).CounterByCodeCorrectHandler).
    Get("/thermometer/:code", 		    (*Context).ThermometerByCodeHandler).
	Put("/thermometer/:code",       	(*Context).ThermometerByCodeAddReadingHandler).
    Get("/flag/:code", 		            (*Context).FlagByCodeHandler).
    Put("/flag/:code",     		    	(*Context).FlagByCodeChangeStateHandler)

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

    c.influxSession = isf()
	c.values = &values
	next(rw, r)
}

func (c *Context) InfluxSessionMiddleware(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
    c.influxSession = isf()
    next(rw, r)
    c.influxSession.Close()
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
		log.WithField("err", err).WithField("s", s).WithField("name", name).Error("Error parsing uint64")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Malformed parameter " + name))
	}

	return
}

func (c *Context) simpleRead(rw web.ResponseWriter, req *web.Request) {
    code := parseStringPathParameter(req, "code")
    newVal, err := c.influxSession.GetLastValue(code)
    if err != nil {
        log.Error(err)
        rw.WriteHeader(500)
    } else {
        rw.WriteHeader(200)
        rw.Write([]byte(fmt.Sprintf("%d", newVal)))
    }
}

func (c *Context) simpleWrite(rw web.ResponseWriter, req *web.Request) {
    code := parseStringPathParameter(req, "code")
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

    err = c.influxSession.AddPoint(code, int64(newReading))
    if err != nil {
        log.Error(err)
        rw.WriteHeader(500)
    } else {
        rw.WriteHeader(200)
        rw.Write([]byte(fmt.Sprintf("%d", newReading)))
    }
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