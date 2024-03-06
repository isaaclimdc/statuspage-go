package statuspage_test

import (
	"context"
	"os"
	"testing"

	statuspage "github.com/andrewwatson/statuspage-go"
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
			PageID:       page,
			Name:         name,
			Body:         body,
			ComponentIDs: components,
			Status:       status,
		}

		result, err := client.Incident.CreateIncident(context.TODO(), page, statuspage.StatusDegraded,incident)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Result ID: %s", result.ID)

	}
}
