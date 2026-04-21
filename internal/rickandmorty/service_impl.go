package rickandmorty

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/elosanz/demo/pkg/circuitbreaker"
)

// rickAndMortyService implements RickAndMortyService.
type rickAndMortyService struct {
	client ExternalRickAndMortyClient
	cb     *circuitbreaker.CircuitBreaker
}

// NewRickAndMortyService constructs a rickAndMortyService with the given external client.
func NewRickAndMortyService(client ExternalRickAndMortyClient) RickAndMortyService {
	cb := circuitbreaker.New(circuitbreaker.Config{
		MaxFailures:       3,
		HalfOpenSuccesses: 1,
		Timeout:           10 * time.Second,
	})
	return &rickAndMortyService{client: client, cb: cb}
}

func (s *rickAndMortyService) GetByID(ctx context.Context, id int64) (*Character, error) {
	var character *Character
	err := s.cb.Execute(ctx, func(ctx context.Context) error {
		var fetchErr error
		character, fetchErr = s.client.FetchByID(ctx, id)
		return fetchErr
	})

	if err != nil {
		if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			return nil, fmt.Errorf("rick and morty api is temporarily unavailable: %w", err)
		}
		return nil, fmt.Errorf("fetching character %d: %w", id, err)
	}
	if character == nil {
		return nil, ErrCharacterNotFound
	}
	return character, nil
}

func (s *rickAndMortyService) GetCharacters(ctx context.Context, page int, name string, status string, species string) (CharacterPage, error) {
	var result CharacterPage
	err := s.cb.Execute(ctx, func(ctx context.Context) error {
		var fetchErr error
		result, fetchErr = s.client.FetchCharacters(ctx, page, name, status, species)
		return fetchErr
	})

	if err != nil {
		if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			return CharacterPage{}, fmt.Errorf("rick and morty api is temporarily unavailable: %w", err)
		}
		return CharacterPage{}, fmt.Errorf("fetching characters: %w", err)
	}
	return result, nil
}
