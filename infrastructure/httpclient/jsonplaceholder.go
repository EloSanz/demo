// Package httpclient provides HTTP adapter implementations for external APIs.
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/elosanz/demo/internal/user"
)

// jsonPlaceholderUserDTO is the wire format from the JSONPlaceholder API.
type jsonPlaceholderUserDTO struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

// JSONPlaceholderClient implements user.ExternalUserClient against JSONPlaceholder.
type JSONPlaceholderClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewJSONPlaceholderClient constructs a JSONPlaceholderClient.
func NewJSONPlaceholderClient(baseURL string, httpClient *http.Client) user.ExternalUserClient {
	return &JSONPlaceholderClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *JSONPlaceholderClient) FetchAll(ctx context.Context) ([]user.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users", nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching all users from jsonplaceholder: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jsonplaceholder responded with %d", resp.StatusCode)
	}

	var dtos []jsonPlaceholderUserDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("decoding jsonplaceholder users: %w", err)
	}

	users := make([]user.User, len(dtos))
	for i, dto := range dtos {
		users[i] = mapJPUserToDomain(dto)
	}
	return users, nil
}

func (c *JSONPlaceholderClient) FetchByID(ctx context.Context, id int64) (*user.User, error) {
	url := c.baseURL + "/users/" + strconv.FormatInt(id, 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request for user %d: %w", id, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching user %d from jsonplaceholder: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jsonplaceholder responded with %d for user %d", resp.StatusCode, id)
	}

	var dto jsonPlaceholderUserDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("decoding jsonplaceholder user %d: %w", id, err)
	}

	u := mapJPUserToDomain(dto)
	return &u, nil
}

// mapJPUserToDomain converts a JSONPlaceholder wire DTO to the domain User type.
func mapJPUserToDomain(dto jsonPlaceholderUserDTO) user.User {
	return user.User{
		ID:      dto.ID,
		Name:    dto.Name,
		Email:   dto.Email,
		Phone:   dto.Phone,
		Website: dto.Website,
	}
}
