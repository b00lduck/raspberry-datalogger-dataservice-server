package orm

type EventType uint8

const (
	TICK EventType = 1 + iota	// Tick of counter
	ABS_CORR					// Absolute correction of reading
	ABS_READ					// Absolute reading
	LAST						// Automatically appendend event with current time
)

type CounterEvent struct {
	ID        	uint 			`gorm:"primary_key"`    // primary key
	CounterID	uint 			`sql:"index"`			// foreign key of the counter
	Timestamp	uint64 			`sql:"index`			// timestamp of the event
	EventType 	EventType								// see above
	Delta		int64 	 								// delta multiplied by 1000
	Reading		uint64  								// reading multiplied by 1000
}

func NewLastCounterEvent(counter Counter) CounterEvent {
	return CounterEvent{
		CounterID: uint(counter.ID),
		Timestamp: GetNow(),
		EventType: LAST,
		Delta:     0,
		Reading:   counter.Reading}
}

func NewTickCounterEvent(counter Counter) CounterEvent {
	return CounterEvent{
		CounterID: uint(counter.ID),
		Timestamp: GetNow(),
		EventType: TICK,
		Delta:     int64(counter.TickAmount),
		Reading:   counter.Reading + counter.TickAmount}
}

func NewAbsCorrCounterEvent(counter Counter, delta int64) CounterEvent {
	return CounterEvent{
		CounterID: uint(counter.ID),
		Timestamp: GetNow(),
		EventType: ABS_CORR,
		Delta:     delta,
		Reading:   uint64(int64(counter.Reading) + delta)}
}