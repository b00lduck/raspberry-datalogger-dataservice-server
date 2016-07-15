package initialization
import (
	"github.com/jinzhu/gorm"
	"github.com/b00lduck/raspberry-datalogger-dataservice-server/orm"
)

type CounterChecker struct {
	db *gorm.DB
}

var counters = []orm.Counter {
	{ Code:	"GAS_1", Name: "Erdgas", Unit: "m³", TickAmount: 100},
	{ Code:	"WAT_1", Name: "Wasser Hauptzähler", Unit: "m³", TickAmount: 10},
	{ Code:	"ELE_1", Name: "Strom", Unit: "kWh", TickAmount: 10}}

var thermometers = []orm.Thermometer {
	{ Code:	"HEIZ_AUSSEN", Name: "Heizung Aussentemperatur"},
	{ Code:	"HEIZ_KESSEL", Name: "Heizung Kesseltemperatur"},
	{ Code:	"HEIZ_BRAUCHW", Name: "Heizung Brauchwassertemperatur"},
}

var flags = []orm.Flag {
	{ Code: "HEIZ_BWZP", Name: "Heizung Brauchwasser Zirkulationspumpe"},
	{ Code: "HEIZ_BWLP", Name: "Heizung Brauchwasser Ladepumpe"},
	{ Code: "HEIZ_WINT", Name: "Heizung Winterbetrieb"},
	{ Code: "HEIZ_UWP", Name: "Heizung Umwälzpumpe"},
	{ Code: "HEIZ_BRENN", Name: "Heizung Brenner"},
	{ Code: "HEIZ_TAG", Name: "Heizung Tagbetrieb (Nachtabsenkung aus)"},
}

func NewCounterChecker(db *gorm.DB) (counterChecker *CounterChecker) {
	counterChecker = new(CounterChecker)
	counterChecker.db = db
	return
}

func (c *CounterChecker) CheckCounters() {
	for i := range counters {
		c.checkCounter(counters[i])
	}
	for i := range thermometers {
		c.checkThermometer(thermometers[i])
	}
	for i := range flags {
		c.checkFlag(flags[i])
	}
}

// Check if the given counter exists in the database and create it if not
func (c *CounterChecker) checkCounter(counter orm.Counter) {
	c.db.Where(orm.Counter{Code: counter.Code}).FirstOrCreate(&counter)
}

// Check if the given thermometer exists in the database and create it if not
func (c *CounterChecker) checkThermometer(thermometer orm.Thermometer) {
	c.db.Where(orm.Thermometer{Code: thermometer.Code}).FirstOrCreate(&thermometer)
}

// Check if the given flag exists in the database and create it if not
func (c *CounterChecker) checkFlag(flag orm.Flag) {
	c.db.Where(orm.Flag{Code: flag.Code}).FirstOrCreate(&flag)
}