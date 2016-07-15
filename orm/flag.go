package orm
import (
	"github.com/jinzhu/gorm"
)

type Flag struct {
	gorm.Model
	Code			string			`sql:"index"`
	Name			string			`sql:"index"`
	State			uint8
	LastChange		uint64			// Timestamp of last change
}