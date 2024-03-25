package statuspage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	statuspage "github.com/andrewwatson/statuspage-go"
	"github.com/stretchr/testify/assert"
)

func TestGetIncident(t *testing.T) {

	if *integration {

		token := os.Getenv("STATUSPAGE_API_TOKEN")
		page := os.Getenv(("STATUSPAGE_API_PAGE"))
		client := statuspage.NewClient(token, nil)

		incident, err := client.Incident.GetIncident(context.TODO(), page, "rgrk1jkj8v8p")
		// t.Logf("%#v", incident)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		t.Logf("Name: %#v", incident)
		// t.Logf("Component: %s", incident.Components[0].Name)

	}
}

func TestClearIncident(t *testing.T) {
	if *integration {
		token := os.Getenv("STATUSPAGE_API_TOKEN")
		page := os.Getenv(("STATUSPAGE_API_PAGE"))
		client := statuspage.NewClient(token, nil)

		existing, err := client.Incident.GetIncident(context.TODO(), page, *incidentID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Status %s", existing.Status)

		existing.Status = statuspage.StatusResolved
		existing.Body = "THIS HAS BEEN CLEARED"
		existing.DeliverNotifications = false

		_, err = client.Incident.UpdateIncident(context.TODO(), page, statuspage.StatusOperational, *existing)
		if err != nil {
			t.Error(err)
		}

	}
}
func TestUpdateIncident(t *testing.T) {
	if *integration {
		token := os.Getenv("STATUSPAGE_API_TOKEN")
		page := os.Getenv(("STATUSPAGE_API_PAGE"))
		client := statuspage.NewClient(token, nil)

		existing, err := client.Incident.GetIncident(context.TODO(), page, *incidentID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Status %s", existing.Status)

		existing.Status = statuspage.StatusMonitoring
		existing.Body = "THIS HAS BEEN MONITORED"
		_, err = client.Incident.UpdateIncident(context.TODO(), page, statuspage.StatusMajorOutage, *existing)
		if err != nil {
			t.Error(err)
		}

	}
}
func TestCreateIncident(t *testing.T) {

	if *integration {

		token := os.Getenv("STATUSPAGE_API_TOKEN")
		page := os.Getenv(("STATUSPAGE_API_PAGE"))
		client := statuspage.NewClient(token, nil)

		components := []string{"qw1nh8v4gxsv"}

		status := statuspage.StatusInvestigating
		name := "Test Incident"
		body := "There is something going on.  We'll figure it out eventually."

		incident := statuspage.Incident{
			PageID:               page,
			Name:                 name,
			Body:                 body,
			ComponentIDs:         components,
			Status:               status,
			DeliverNotifications: false,
		}

		result, err := client.Incident.CreateIncident(context.TODO(), page, statuspage.StatusDegraded, incident)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Result ID: %s", result.ID)

	}
}

func TestJSONEncodeIncident(t *testing.T) {

	components := []string{"qw1nh8v4gxsv"}
	page := os.Getenv(("STATUSPAGE_API_PAGE"))
	status := statuspage.StatusInvestigating
	name := "Test Incident"
	body := "There is something going on.  We'll figure it out eventually."

	incident := statuspage.Incident{
		PageID:               page,
		Name:                 name,
		Body:                 body,
		ComponentIDs:         components,
		Status:               status,
		DeliverNotifications: false,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(incident)

	// encoded, err := json.Marshal(incident)
	if err != nil {
		t.Error(err)
		return
	}

	encoded := buf.String()

	assert.NotEmpty(t, encoded)
	t.Logf("ENCODED: %s", encoded)
}
