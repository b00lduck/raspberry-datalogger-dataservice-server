package rest

import (
	"github.com/gocraft/web"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/orm"
	"net/http"
	"io/ioutil"
	"strconv"
)

// Get all flags
func (c *Context) FlagHandler(rw web.ResponseWriter, req *web.Request) {
	var flags []orm.Counter
	db.Find(&flags)
	marshal(rw, flags)
}

// Get specific flag by code
func (c *Context) FlagByCodeHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var flag orm.Flag
	db.Where(&orm.Flag{Code: code}).First(&flag)

	if (flag.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Flag not found"))
		return
	}

	marshal(rw, flag)
}

// Chenge flag state by code
func (c *Context) FlagByCodeChangeStateHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var flag orm.Flag
	db.Where(&orm.Flag{Code: code}).First(&flag)

	if (flag.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Flag not found"))
		return
	}

	body, err := ioutil.ReadAll(req.Body);
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not read body"))
		return
	}

	state, err := strconv.Atoi(string(body))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Could not parse state"))
		return
	}

	flagState := orm.NewFlagState(flag, uint8(state))
	db.Create(&flagState)

	flag.State = uint8(state)
	flag.LastChange = flagState.Timestamp
	db.Save(flag)

	marshal(rw, flagState)
}

// Get flag states in a optionally given time range
// Query parameters: start,end
func (c *Context) FlagByCodeGetStatesHandler(rw web.ResponseWriter, req *web.Request) {

	code := parseStringPathParameter(req, "code")
	var flag orm.Flag
	db.Where(&orm.Flag{Code: code}).First(&flag)

	if (flag.ID == 0) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Flag not found"))
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

	var flagStates []orm.FlagState
	orm.GetOrderedWindowedQuery(db, "flag_id", flag.ID, start, end).Find(&flagStates)

	var startReading orm.FlagState
	var endReading orm.FlagState

	if len(flagStates) == 0 {

		startReading = orm.FlagState{
			State: flag.State,
			Timestamp: start,
			FlagID: flag.ID,
		}

		endReading = orm.FlagState{
			State: flag.State,
			Timestamp: end,
			FlagID: flag.ID,
		}

	} else {

		if db.Where("timestamp < ? and flag_id = ?", start, flag.ID).Order("timestamp desc").First(&startReading).RecordNotFound() {
			startReading = orm.FlagState{
				State: flagStates[0].State,
				Timestamp: start,
				FlagID: flag.ID,
			}
		} else {
			startReading.Timestamp = start
		}

		if db.Where("timestamp > ? and flag_id = ?", end, flag.ID).Order("timestamp asc").First(&endReading).RecordNotFound() {
			endReading = orm.FlagState{
				State: flagStates[len(flagStates) - 1].State,
				Timestamp: end,
				FlagID: flag.ID,
			}
		} else {
			endReading.Timestamp = end
		}

	}

	flagStates = append([]orm.FlagState{startReading}, flagStates...)
	flagStates = append(flagStates, endReading)

	marshal(rw, flagStates)

}
