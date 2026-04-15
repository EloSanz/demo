package rickandmorty

import (
	"context"
	"fmt"
)

// rickAndMortyService implements RickAndMortyService.
type rickAndMortyService struct {
	client ExternalRickAndMortyClient
}

// NewRickAndMortyService constructs a rickAndMortyService with the given external client.
func NewRickAndMortyService(client ExternalRickAndMortyClient) RickAndMortyService {
	return &rickAndMortyService{client: client}
}

func (s *rickAndMortyService) GetByID(ctx context.Context, id int64) (*Character, error) {
	character, err := s.client.FetchByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetching character %d: %w", id, err)
	}
	if character == nil {
		return nil, ErrCharacterNotFound
	}
	return character, nil
}

func (s *rickAndMortyService) GetCharacters(ctx context.Context, page int, name string, status string, species string) (CharacterPage, error) {
	result, err := s.client.FetchCharacters(ctx, page, name, status, species)
	if err != nil {
		return CharacterPage{}, fmt.Errorf("fetching characters: %w", err)
	}
	return result, nil
}
