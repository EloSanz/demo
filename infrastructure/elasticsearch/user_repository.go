package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elosanz/demo/internal/user"
)

type userSearchRepository struct {
	client *elasticsearch.Client
	index  string
}

func NewUserSearchRepository(client *elasticsearch.Client) user.UserSearchRepository {
	return &userSearchRepository{
		client: client,
		index:  "users",
	}
}

func (r *userSearchRepository) Index(ctx context.Context, u user.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshaling user for elasticsearch: %w", err)
	}

	res, err := r.client.Index(
		r.index,
		bytes.NewReader(data),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(fmt.Sprintf("%d", u.ID)),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch index request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	slog.Debug("user indexed in elasticsearch", "id", u.ID)
	return nil
}

func (r *userSearchRepository) Search(ctx context.Context, query string) ([]user.User, error) {
	var buf bytes.Buffer
	searchQuery := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":     query,
				"fields":    []string{"name", "email", "phone"},
				"fuzziness": "AUTO",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	var rMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&rMap); err != nil {
		return nil, err
	}

	var users []user.User
	for _, hit := range rMap["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		var u user.User
		b, _ := json.Marshal(source)
		json.Unmarshal(b, &u)
		users = append(users, u)
	}

	return users, nil
}
