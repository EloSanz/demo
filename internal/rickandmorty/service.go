package rickandmorty

import "context"

type RickAndMortyService interface {
	GetByID(ctx context.Context, id int64) (*Character, error)
	GetCharacters(ctx context.Context, page int, name string, status string, species string) (CharacterPage, error)
}

// ExternalRickAndMortyClient is the port for the public Rick and Morty API.
type ExternalRickAndMortyClient interface {
	FetchByID(ctx context.Context, id int64) (*Character, error)
	FetchCharacters(ctx context.Context, page int, name string, status string, species string) (CharacterPage, error)
}
