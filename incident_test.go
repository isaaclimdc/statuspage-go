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

		incident, err := client.Incident.GetIncident(context.TODO(), page, "jmnlg1ckt15p")
		// t.Logf("%#v", incident)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		t.Logf("Name: %s", *incident.Name)
		t.Logf("Component: %s", *incident.Components[0].Name)

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
			PageID:       &page,
			Name:         &name,
			Body:         &body,
			ComponentIDs: components,
			Status:       &status,
		}

		_, err := client.Incident.CreateIncident(context.TODO(), page, incident)
		if err != nil {
			t.Error(err)
		}
	}
}
