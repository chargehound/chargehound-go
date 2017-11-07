package chargehound_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chargehound/chargehound-go"
)

func TestRetrieveDispute(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx" {
			t.Error("Incorrect path.")
		}

		if r.Header.Get("User-Agent") != "Chargehound/v1 GoBindings/"+ch.Version {
			t.Error("Incorrect version.")
		}

		if r.Header.Get("Authorization") != "Basic YXBpX2tleTo=" {
			t.Error("Incorrect authorization.")
		}

		if r.Header.Get("Chargehound-Version") != "" {
			t.Error("Incorrect version.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Retrieve(&chargehound.RetrieveDisputeParams{ID: "dp_xxx"})
	if err != nil {
		t.Error(err)
	}
}

func TestRetrieveDisputeResponse(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx/response" {
			t.Error("Incorrect path.")
		}

		if r.Header.Get("User-Agent") != "Chargehound/v1 GoBindings/"+ch.Version {
			t.Error("Incorrect version.")
		}

		if r.Header.Get("Authorization") != "Basic YXBpX2tleTo=" {
			t.Error("Incorrect authorization.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Response(&chargehound.RetrieveDisputeParams{ID: "dp_xxx"})
	if err != nil {
		t.Error(err)
	}
}

type TestTransport struct {
	Called    bool
	Transport *http.Transport
}

func (tt *TestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	tt.Called = true
	return tt.Transport.RoundTrip(r)
}

func TestOptHTTPClientDispute(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx" {
			t.Error("Incorrect path.")
		}

		if r.Header.Get("User-Agent") != "Chargehound/v1 GoBindings/"+ch.Version {
			t.Error("Incorrect version.")
		}

		if r.Header.Get("Authorization") != "Basic YXBpX2tleTo=" {
			t.Error("Incorrect authorization.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	tt := TestTransport{Transport: &http.Transport{}}

	_, err = ch.Disputes.Retrieve(&chargehound.RetrieveDisputeParams{
		ID:            "dp_xxx",
		OptHTTPClient: &http.Client{Transport: &tt},
	})

	if err != nil {
		t.Error(err)
	}

	if tt.Called == false {
		t.Error("Optional client not used.")
	}
}

func TestListDisputes(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes" {
			t.Error("Incorrect path.")
		}

		if r.URL.RawQuery != "starting_after=dp_yyy" {
			t.Error("Incorrect query.")
		}

		if r.Header.Get("User-Agent") != "Chargehound/v1 GoBindings/"+ch.Version {
			t.Error("Incorrect version.")
		}

		if r.Header.Get("Authorization") != "Basic YXBpX2tleTo=" {
			t.Error("Incorrect authorization.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.List(&chargehound.ListDisputesParams{StartingAfter: "dp_yyy"})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDisputeFields(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["template"] != "tmpl_1" {
			t.Error("Incorrect template id.")
		}

		fields := b["fields"].(map[string]interface{})

		if fields["f1"] != "v1" {
			t.Error("Incorrect string field.")
		}

		if fields["f2"] != 2.0 {
			t.Error("Incorrect int field.")
		}

		_, ok := b["products"]

		if ok {
			t.Error("Incorrect products data.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Update(&chargehound.UpdateDisputeParams{
		ID:       "dp_xxx",
		Template: "tmpl_1",
		Fields: map[string]interface{}{
			"f1": "v1",
			"f2": 2,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDisputeProducts(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["template"] != "tmpl_1" {
			t.Error("Incorrect template id.")
		}

		products := b["products"].([]interface{})
		product1 := products[0].(map[string]interface{})

		if product1["name"] != "prod1" {
			t.Error("Incorrect product name.")
		}

		if product1["amount"] != 100.0 {
			t.Error("Incorrect product amount.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Update(&chargehound.UpdateDisputeParams{
		ID:       "dp_xxx",
		Template: "tmpl_1",
		Products: []chargehound.Product{{Name: "prod1", Amount: 100}},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDisputeUserID(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["user_id"] != "acct_xxx" {
			t.Error("Incorrect account id.")
		}

		if _, ok := b["force"]; ok {
			t.Error("Sending force param erroneously.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Update(&chargehound.UpdateDisputeParams{
		ID:     "dp_xxx",
		UserID: "acct_xxx",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDisputeForce(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx/submit" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["force"] != true {
			t.Error("Incorrect force parameter.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Submit(&chargehound.UpdateDisputeParams{
		ID:    "dp_xxx",
		Force: true,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDisputeCharge(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx/submit" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["charge"] != "ch_XXX" {
			t.Error("Incorrect charge parameter.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Submit(&chargehound.UpdateDisputeParams{
		ID:     "dp_xxx",
		Charge: "ch_XXX",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestQueueDispute(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx/submit" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["queue"] != true {
			t.Error("Incorrect queue parameter.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Submit(&chargehound.UpdateDisputeParams{
		ID:     "dp_xxx",
		Charge: "ch_XXX",
		Queue:  true,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestAcceptDispute(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx/accept" {
			t.Error("Incorrect path.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Accept(&chargehound.AcceptDisputeParams{
		ID: "dp_xxx",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestResponseCode(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	list, err := ch.Disputes.List(&chargehound.ListDisputesParams{StartingAfter: "dp_yyy"})
	if err != nil {
		t.Error(err)
	}
	if list.Response.Status != 200 {
		t.Error("Missing response status code.")
	}
}

func TestCreateDispute(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes" {
			t.Error("Incorrect path.")
		}

		decoder := json.NewDecoder(r.Body)
		b := make(map[string]interface{})
		err := decoder.Decode(&b)
		if err != nil {
			t.Error(err)
		}

		if b["id"] != "dp_xxx" {
			t.Error("Incorrect dispute id.")
		}

		if b["template"] != "tmpl_1" {
			t.Error("Incorrect template id.")
		}

		fields := b["fields"].(map[string]interface{})

		if fields["f1"] != "v1" {
			t.Error("Incorrect string field.")
		}

		if fields["f2"] != 2.0 {
			t.Error("Incorrect int field.")
		}

		_, ok := b["products"]

		if ok {
			t.Error("Incorrect products data.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Create(&chargehound.CreateDisputeParams{
		ID:       "dp_xxx",
		Template: "tmpl_1",
		Fields: map[string]interface{}{
			"f1": "v1",
			"f2": 2,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestOverrideAPIVersion(t *testing.T) {
	ch := chargehound.New("api_key",
		&chargehound.ClientParams{APIVersion: "1999-01-01"})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Error("Incorrect method.")
		}

		if r.URL.Path != "/v1/disputes/dp_xxx" {
			t.Error("Incorrect path.")
		}

		if r.Header.Get("User-Agent") != "Chargehound/v1 GoBindings/"+ch.Version {
			t.Error("Incorrect version.")
		}

		if r.Header.Get("Authorization") != "Basic YXBpX2tleTo=" {
			t.Error("Incorrect authorization.")
		}

		if r.Header.Get("Chargehound-Version") != "1999-01-01" {
			t.Error("Incorrect version.")
		}

		json.NewEncoder(w).Encode(chargehound.Dispute{ID: "dp_xxx"})
	}))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch.Host = url.Host
	ch.Protocol = url.Scheme + "://"

	_, err = ch.Disputes.Retrieve(&chargehound.RetrieveDisputeParams{ID: "dp_xxx"})
	if err != nil {
		t.Error(err)
	}
}
