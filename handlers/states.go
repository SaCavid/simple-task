package handlers

import "fmt"

type SourceType int

const (
	game SourceType = iota
	server
	payment
)

var (
	// must be unique names
	// index must be same in constants
	SourceTypes = [...]string{"game", "server", "payment"}
)

func (s SourceType) String() string {
	return SourceTypes[s]
}

func (s SourceType) IndexOf(name string) (int, error) {

	for k, v := range SourceTypes {
		if v == name {
			return k, nil
		}
	}

	return -1, fmt.Errorf("not acceptable source type")
}
