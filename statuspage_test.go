package statuspage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	statuspage "github.com/isaaclimdc/statuspage-go"
)

var referenceTime = time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)

var integration = flag.Bool("integration", false, "enable integration testing")

var incidentID = flag.String("incident", "", "for update/clear")

// var offline = flag.Bool("offline", true, "run offline tests")

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/test"
)

// setup sets up a test HTTP server along with a statuspage.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *statuspage.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Statuspage client being tested and is
	// configured to use test server.
	client = statuspage.NewClient("test-token", nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

// Helper function to test that a value is marshalled to JSON as expected.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %v", v)
	}

	w := new(bytes.Buffer)
	err = json.Compact(w, []byte(want))
	if err != nil {
		t.Errorf("String is not valid json: %s", want)
	}

	if w.String() != string(j) {
		t.Errorf("json.Marshal(%q) returned %s, want %s", v, j, w)
	}

	// now go the other direction and make sure things unmarshal as expected
	u := reflect.ValueOf(v).Interface()
	if err := json.Unmarshal([]byte(want), u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v", want)
	}

	if !reflect.DeepEqual(v, u) {
		t.Errorf("json.Unmarshal(%q) returned %s, want %s", want, u, v)
	}
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int32 is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it.
func Int32(v int32) *int32 { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// func TestUpdateIncident(t *testing.T) {
// 	if *integration {
// 		t.Log("INTEGRATION")

// 		token := os.Getenv("STATUSPAGE_API_TOKEN")
// 		page := os.Getenv(("STATUSPAGE_API_PAGE"))
// 		client := statuspage.NewClient(token, nil)

// 		incident := statuspage.Incident{
// 			Name:                 "Integration Test",
// 			Body:                 "Something is broken",
// 			DeliverNotifications: true,
// 			Status:               statuspage.StatusIdentified,
// 			ComponentIDs: []string{"qw1nh8v4gxsv"},
// 		}

// 		result, err := client.Incident.CreateIncident(context.TODO(), page, incident)
// 		if err != nil {
// 			t.Errorf("Create Incident %s", err.Error())
// 			return
// 		}

// 		updatedIncident := statuspage.Incident{
// 			ID: result.ID,
// 			Name:                 incident.Name,
// 			Body:                 "Issue has been resolved",
// 			DeliverNotifications: false,
// 			Status:               statuspage.StatusResolved,
// 			ComponentIDs: []string{"qw1nh8v4gxsv"},

// 		}

// 		time.Sleep(10 * time.Second)

// 		_, err = client.Incident.UpdateIncident(context.TODO(), page, updatedIncident)
// 		if err != nil {
// 			t.Errorf("Update Incident %s", err.Error())
// 			return
// 		}
// 	}
// }

func TestGetGroups(t *testing.T) {
	if *integration {
		t.Log("INTEGRATION")

		token := os.Getenv("STATUSPAGE_API_TOKEN")
		page := os.Getenv(("STATUSPAGE_API_PAGE"))
		client := statuspage.NewClient(token, nil)

		groups, err := client.GetAllGroupsAndComponents(context.TODO(), page)
		if err != nil {
			t.Error(err)
		}

		if len(groups) < 1 {
			t.Fail()
		}
	}

}
func TestIntegration(t *testing.T) {
	if *integration {
		t.Log("INTEGRATION")

		token := os.Getenv("STATUSPAGE_API_TOKEN")

		page := os.Getenv(("STATUSPAGE_API_PAGE"))

		client := statuspage.NewClient(token, nil)

		groups, err := client.Group.GetGroups(context.TODO(), page)
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}

		t.Logf("groups: %d: %s", len(groups), groups[0].ID)

		pGroup, err := client.Group.GetGroup(context.TODO(), page, groups[0].ID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Components: %d: %s", len(pGroup.Components), pGroup.Components[0])

		components, err := client.GetComponentsFromGroup(context.TODO(), page, groups[0].ID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Components: %d", len(components))
	}
}
