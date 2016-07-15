package orm

import (
	"time"
	"github.com/jinzhu/gorm"
)

func GetNow() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func GetOrderedWindowedQuery(db *gorm.DB, idfield string, id uint, start uint64, end uint64) *gorm.DB {

	var q *gorm.DB

	switch {
	case start > 0 && end > 0:
		q = db.Where(idfield + " = ? AND timestamp >= ? AND timestamp <= ?", id, start, end)
	case end > 0:
		q = db.Where(idfield + " = ? AND timestamp <= ?", id, end)
	case start > 0:
		q = db.Where(idfield + " = ? AND timestamp >= ?", id, start)
	default:
		q = db.Where(idfield + " = ?", id)
	}

	return q.Order("timestamp")

}
