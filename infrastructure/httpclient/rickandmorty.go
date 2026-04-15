package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/elosanz/demo/internal/rickandmorty"
)

// rmCharacterDTO is the wire format from the Rick and Morty API for a single character.
type rmCharacterDTO struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
	Image   string `json:"image"`
	URL     string `json:"url"`
	Created string `json:"created"`
}

// rmPageInfoDTO is the pagination metadata from the Rick and Morty API.
type rmPageInfoDTO struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

// rmPageDTO is the paginated response from the Rick and Morty API.
type rmPageDTO struct {
	Info    rmPageInfoDTO    `json:"info"`
	Results []rmCharacterDTO `json:"results"`
}

// RickAndMortyHTTPClient implements rickandmorty.ExternalRickAndMortyClient.
type RickAndMortyHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewRickAndMortyHTTPClient constructs a RickAndMortyHTTPClient.
func NewRickAndMortyHTTPClient(baseURL string, httpClient *http.Client) rickandmorty.ExternalRickAndMortyClient {
	return &RickAndMortyHTTPClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *RickAndMortyHTTPClient) FetchByID(ctx context.Context, id int64) (*rickandmorty.Character, error) {
	endpoint := c.baseURL + "/character/" + strconv.FormatInt(id, 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for character %d: %w", id, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching character %d: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rick and morty api responded with %d for character %d", resp.StatusCode, id)
	}

	var dto rmCharacterDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("decoding character %d: %w", id, err)
	}

	ch := mapRMCharacterToDomain(dto)
	return &ch, nil
}

func (c *RickAndMortyHTTPClient) FetchCharacters(ctx context.Context, page int, name string, status string, species string) (rickandmorty.CharacterPage, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if name != "" {
		params.Set("name", name)
	}
	if status != "" {
		params.Set("status", status)
	}
	if species != "" {
		params.Set("species", species)
	}

	endpoint := c.baseURL + "/character/?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return rickandmorty.CharacterPage{}, fmt.Errorf("building characters request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return rickandmorty.CharacterPage{}, fmt.Errorf("fetching characters: %w", err)
	}
	defer resp.Body.Close()

	// 404 from the Rick and Morty API means no results for the given filters.
	if resp.StatusCode == http.StatusNotFound {
		return rickandmorty.CharacterPage{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return rickandmorty.CharacterPage{}, fmt.Errorf("rick and morty api responded with %d", resp.StatusCode)
	}

	var dto rmPageDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return rickandmorty.CharacterPage{}, fmt.Errorf("decoding characters page: %w", err)
	}

	return mapRMPageToDomain(dto), nil
}

// mapRMCharacterToDomain converts a Rick and Morty wire DTO to the domain Character type.
func mapRMCharacterToDomain(dto rmCharacterDTO) rickandmorty.Character {
	return rickandmorty.Character{
		ID:      dto.ID,
		Name:    dto.Name,
		Status:  dto.Status,
		Species: dto.Species,
		Type:    dto.Type,
		Gender:  dto.Gender,
		Image:   dto.Image,
		URL:     dto.URL,
		Created: dto.Created,
	}
}

// mapRMPageToDomain converts a Rick and Morty page wire DTO to the domain CharacterPage type.
func mapRMPageToDomain(dto rmPageDTO) rickandmorty.CharacterPage {
	results := make([]rickandmorty.Character, len(dto.Results))
	for i, ch := range dto.Results {
		results[i] = mapRMCharacterToDomain(ch)
	}
	return rickandmorty.CharacterPage{
		Info: rickandmorty.PageInfo{
			Count: dto.Info.Count,
			Pages: dto.Info.Pages,
			Next:  dto.Info.Next,
			Prev:  dto.Info.Prev,
		},
		Results: results,
	}
}
