package handlers

import "fmt"

type SourceType int

// const for source type of requests
const (
	game SourceType = iota
	server
	payment
)

var (
	// must be unique names
	// index must be same as in constants
	SourceTypes = [...]string{"game", "server", "payment"}
)

// get source type as string
func (s SourceType) String() string {
	return SourceTypes[s]
}

// get index of source type
func (s SourceType) IndexOf(name string) (int, error) {

	for k, v := range SourceTypes {
		if v == name {
			return k, nil
		}
	}

	return -1, fmt.Errorf("not acceptable source type")
}
