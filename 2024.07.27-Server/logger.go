package main

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type EventType byte

type TransactionLogger interface {
	WriteDelete(id uint64)
	WritePut(id uint64, title, artist, price string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type Event struct {
	EventType EventType // The action token
	Id        uint64    // The id affected by this transaction
	Title     string    // The value of a PUT the transaction
	Artist    string    //
	Prise     string    //
}
