package statuspage

import (
	"context"
)

type GroupService service

type Group struct {
	ID             string      `json:"id,omitempty"`
	PageID         string      `json:"page_id,omitempty"`
	CreatedAt      Timestamp   `json:"created_at,omitempty"`
	UpdatedAt      Timestamp   `json:"updated_at,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	Position       int32       `json:"position,omitempty"`
	Components     []string    `json:"components,omitempty"`
	FullComponents []Component `json:"fullComponents,omitempty"`
}

// GetGroup returns component group information for a given page and component group id
func (s *GroupService) GetGroup(ctx context.Context, pageID string, groupID string) (*Group, error) {
	path := "v1/pages/" + pageID + "/component-groups/" + groupID
	req, err := s.client.newRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	var group Group
	_, err = s.client.do(ctx, req, &group)

	return &group, err
}

// ListComponents returns a list of all components for a given page id
func (s *GroupService) GetGroups(ctx context.Context, pageID string) ([]Group, error) {
	path := "v1/pages/" + pageID + "/component-groups"
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var groups []Group
	_, err = s.client.do(ctx, req, &groups)

	return groups, err
}
