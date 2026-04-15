package handler

import (
	"github.com/elosanz/demo/internal/rickandmorty"
)

// CharacterResponse is the API representation of a Rick and Morty character.
type CharacterResponse struct {
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

// PageInfoResponse holds pagination metadata.
type PageInfoResponse struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

// CharacterPageResponse is the paginated response for character searches.
type CharacterPageResponse struct {
	Info    PageInfoResponse    `json:"info"`
	Results []CharacterResponse `json:"results"`
}

func mapCharacterToResponse(c rickandmorty.Character) CharacterResponse {
	return CharacterResponse{
		ID:      c.ID,
		Name:    c.Name,
		Status:  c.Status,
		Species: c.Species,
		Type:    c.Type,
		Gender:  c.Gender,
		Image:   c.Image,
		URL:     c.URL,
		Created: c.Created,
	}
}

func mapPageToResponse(p rickandmorty.CharacterPage) CharacterPageResponse {
	results := make([]CharacterResponse, len(p.Results))
	for i, c := range p.Results {
		results[i] = mapCharacterToResponse(c)
	}
	return CharacterPageResponse{
		Info: PageInfoResponse{
			Count: p.Info.Count,
			Pages: p.Info.Pages,
			Next:  p.Info.Next,
			Prev:  p.Info.Prev,
		},
		Results: results,
	}
}
