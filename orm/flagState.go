package orm

type FlagState struct {
	ID        		uint 			`gorm:"primary_key"`    // primary key
	FlagID			uint 			`sql:"index"`			// foreign key of the flag
	Timestamp		uint64 			`sql:"index`			// timestamp of the event
	State			uint8
}

func NewFlagState(flag Flag, state uint8) FlagState {
	return FlagState{
		FlagID: uint(flag.ID),
		Timestamp: GetNow(),
		State: state}
}