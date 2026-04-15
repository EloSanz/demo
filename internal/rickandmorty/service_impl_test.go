package rickandmorty_test

import (
	"context"
	"testing"

	"github.com/elosanz/demo/internal/rickandmorty"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

// ─── Fakes ───────────────────────────────────────────────────────────────────

type fakeRickAndMortyClient struct {
	byID       map[int64]*rickandmorty.Character
	pageResult rickandmorty.CharacterPage
}

func newFakeRickAndMortyClient() *fakeRickAndMortyClient {
	return &fakeRickAndMortyClient{
		byID: make(map[int64]*rickandmorty.Character),
	}
}

func (f *fakeRickAndMortyClient) FetchByID(_ context.Context, id int64) (*rickandmorty.Character, error) {
	return f.byID[id], nil
}

func (f *fakeRickAndMortyClient) FetchCharacters(_ context.Context, _ int, _ string, _ string, _ string) (rickandmorty.CharacterPage, error) {
	return f.pageResult, nil
}

// ─── Factory ─────────────────────────────────────────────────────────────────

type rickAndMortyServiceBuilder struct {
	client *fakeRickAndMortyClient
}

func newRickAndMortyServiceFactory() *rickAndMortyServiceBuilder {
	return &rickAndMortyServiceBuilder{
		client: newFakeRickAndMortyClient(),
	}
}

func (b *rickAndMortyServiceBuilder) withCharacter(id int64, c rickandmorty.Character) *rickAndMortyServiceBuilder {
	b.client.byID[id] = &c
	return b
}

func (b *rickAndMortyServiceBuilder) withPageResult(page rickandmorty.CharacterPage) *rickAndMortyServiceBuilder {
	b.client.pageResult = page
	return b
}

func (b *rickAndMortyServiceBuilder) build() rickandmorty.RickAndMortyService {
	return rickandmorty.NewRickAndMortyService(b.client)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestRickAndMortyService_GetByID_Success(t *testing.T) {
	// Given
	svc := newRickAndMortyServiceFactory().
		withCharacter(1, rickandmorty.Character{ID: 1, Name: "Rick Sanchez", Status: "Alive"}).
		build()

	// When
	character, err := svc.GetByID(ctx, 1)

	// Then
	require.NoError(t, err)
	require.NotNil(t, character)
	require.Equal(t, "Rick Sanchez", character.Name)
}

func TestRickAndMortyService_GetByID_Error_WithCharacterNotFound(t *testing.T) {
	// Given
	svc := newRickAndMortyServiceFactory().build()

	// When
	_, err := svc.GetByID(ctx, 999)

	// Then
	require.ErrorIs(t, err, rickandmorty.ErrCharacterNotFound)
}

func TestRickAndMortyService_GetCharacters_Success(t *testing.T) {
	// Given
	page := rickandmorty.CharacterPage{
		Info:    rickandmorty.PageInfo{Count: 1, Pages: 1},
		Results: []rickandmorty.Character{{ID: 1, Name: "Rick Sanchez"}},
	}
	svc := newRickAndMortyServiceFactory().
		withPageResult(page).
		build()

	// When
	result, err := svc.GetCharacters(ctx, 1, "Rick", "", "")

	// Then
	require.NoError(t, err)
	require.Equal(t, 1, len(result.Results))
	require.Equal(t, "Rick Sanchez", result.Results[0].Name)
}

func TestRickAndMortyService_GetCharacters_Success_WithEmptyResult(t *testing.T) {
	// Given
	svc := newRickAndMortyServiceFactory().build()

	// When
	result, err := svc.GetCharacters(ctx, 1, "nonexistent-character", "", "")

	// Then
	require.NoError(t, err)
	require.Empty(t, result.Results)
}
