package orm

type ThermometerReading struct {
	ID        		uint 			`gorm:"primary_key"`    // primary key
	ThermometerID	uint 			`sql:"index"`			// foreign key of the thermometer
	Timestamp		uint64 			`sql:"index`			// timestamp of the event
	Reading			uint64  								// reading multiplied by 1000
}

func NewThermometerReading(thermometer Thermometer, reading uint64) ThermometerReading {
	return ThermometerReading{
		ThermometerID: uint(thermometer.ID),
		Timestamp: GetNow(),
		Reading:   reading}
}