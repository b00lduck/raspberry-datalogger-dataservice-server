package orm
import (
	"github.com/jinzhu/gorm"
)

type Thermometer struct {
	gorm.Model
	Code			string			`sql:"index"`
	Name			string			`sql:"index"`
	Reading			uint64			// Reading multiplied by 1000
	LastReading		uint64			// Timestamp of last tick
}