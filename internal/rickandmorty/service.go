package rickandmorty

import "context"

// RickAndMortyService defines the use-case contract for the Rick and Morty bounded context.
type RickAndMortyService interface {
	// GetByID returns ErrCharacterNotFound when the character does not exist.
	GetByID(ctx context.Context, id int64) (*Character, error)
	// GetCharacters returns an empty CharacterPage when no characters match the filters.
	GetCharacters(ctx context.Context, page int, name string, status string, species string) (CharacterPage, error)
}

// ExternalRickAndMortyClient is the port for the public Rick and Morty API.
type ExternalRickAndMortyClient interface {
	// FetchByID returns nil when the character is not found (404).
	FetchByID(ctx context.Context, id int64) (*Character, error)
	FetchCharacters(ctx context.Context, page int, name string, status string, species string) (CharacterPage, error)
}
