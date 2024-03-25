package statuspage

import (
	"context"
)

const (
	StatusInvestigating = "investigating"
	StatusIdentified    = "identified"
	StatusMonitoring    = "monitoring"
	StatusResolved      = "resolved"

	StatusOperational   = "operational"
	StatusMaintenance   = "under_maintenance"
	StatusDegraded      = "degraded_performance"
	StatusPartialOutage = "partial_outage"
	StatusMajorOutage   = "major_outage"
)

type IncidentService service

// UpdateComponentRequestBody is the update component request body representation
type UpdateIncidentRequestBody struct {
	Incident IncidentUpdate `json:"incident"`
}

type CreateIncidentRequestBody struct {
	Incident Incident `json:"incident"`
}

type Incident struct {
	ID                   string      `json:"id,omitempty"`
	PageID               string      `json:"page_id,omitempty"`
	CreatedAt            Timestamp   `json:"created_at,omitempty"`
	UpdatedAt            Timestamp   `json:"updated_at,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Body                 string      `json:"body"`
	Status               string      `json:"status,omitempty"`
	Components           []Component `json:"components,omitempty"`
	ComponentIDs         []string    `json:"component_ids,omitempty"`
	DeliverNotifications bool        `json:"deliver_notifications"`
}

type IncidentUpdate struct {
	ID                   string            `json:"id,omitempty"`
	PageID               string            `json:"page_id,omitempty"`
	CreatedAt            Timestamp         `json:"created_at,omitempty"`
	UpdatedAt            Timestamp         `json:"updated_at,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Body                 string            `json:"body"`
	Status               string            `json:"status,omitempty"`
	Components           map[string]string `json:"components,omitempty"`
	ComponentIDs         []string          `json:"component_ids,omitempty"`
	DeliverNotifications bool              `json:"deliver_notifications"`
}

// CreateIncident creates a new incident
func (s *IncidentService) CreateIncident(ctx context.Context, pageID, status string, incident Incident) (*Incident, error) {

	if pageID == "" {
		pageID = s.client.defaultPage
	}

	componentMap := make(map[string]string, 0)

	for _, c := range incident.ComponentIDs {
		componentMap[c] = status
	}

	updateBody := IncidentUpdate{
		ID:                   incident.ID,
		Name:                 incident.Name,
		Body:                 incident.Body,
		Components:           componentMap,
		ComponentIDs:         incident.ComponentIDs,
		Status:               incident.Status,
		DeliverNotifications: incident.DeliverNotifications,
	}

	path := "v1/pages/" + pageID + "/incidents/"
	requestBody := UpdateIncidentRequestBody{updateBody}

	req, err := s.client.newRequest("POST", path, requestBody)

	if err != nil {
		return nil, err
	}

	_, err = s.client.do(ctx, req, &incident)

	return &incident, err
}

// GetGroup returns component group information for a given page and component group id
func (s *IncidentService) GetIncident(ctx context.Context, pageID string, incidentID string) (*Incident, error) {
	path := "v1/pages/" + pageID + "/incidents/" + incidentID
	req, err := s.client.newRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	var incident Incident
	_, err = s.client.do(ctx, req, &incident)

	return &incident, err
}

// UpdateIncident updates a component for a given page and component id
func (s *IncidentService) UpdateIncident(ctx context.Context, pageID, status string, incident Incident) (*Incident, error) {

	if pageID == "" {
		pageID = s.client.defaultPage
	}

	path := "v1/pages/" + pageID + "/incidents/" + incident.ID

	componentMap := make(map[string]string, 0)

	for _, c := range incident.Components {
		componentMap[c.ID] = status
	}

	updateBody := IncidentUpdate{
		ID:                   incident.ID,
		Name:                 incident.Name,
		Body:                 incident.Body,
		Components:           componentMap,
		ComponentIDs:         incident.ComponentIDs,
		Status:               incident.Status,
		DeliverNotifications: incident.DeliverNotifications,
	}

	payload := UpdateIncidentRequestBody{Incident: updateBody}
	req, err := s.client.newRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	var updatedIncident Incident
	_, err = s.client.do(ctx, req, &updatedIncident)

	return &updatedIncident, err
}
