package rickandmorty

import "errors"

// Character is the domain entity for a Rick and Morty character.
type Character struct {
	ID      int64
	Name    string
	Status  string
	Species string
	Type    string
	Gender  string
	Image   string
	URL     string
	Created string
}

// PageInfo holds pagination metadata from the Rick and Morty API.
type PageInfo struct {
	Count int
	Pages int
	Next  string
	Prev  string
}

// CharacterPage is a paginated result set.
type CharacterPage struct {
	Info    PageInfo
	Results []Character
}

// ErrCharacterNotFound is returned when a character cannot be found by the given ID.
var ErrCharacterNotFound = errors.New("character not found")
