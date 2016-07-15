package orm
import (
	"github.com/jinzhu/gorm"
)

type Counter struct {
	gorm.Model
	Code			string			`sql:"index"`
	Name			string			`sql:"index"`
	Unit 			string
	TickAmount		uint64			// Tick amount multiplied by 1000
	Reading			uint64			// Reading multiplied by 1000
	LastTick		uint64			// Timestamp of last tick
}