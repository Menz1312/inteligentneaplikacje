package dbase

type rowscan interface {
	Scan(...any) error
}
