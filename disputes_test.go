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
	ch := chargehound.New("api_key")

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

	_, err = ch.Disputes.Retrieve(&chargehound.RetrieveDisputeParams{ID: "dp_xxx"})
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
	ch := chargehound.New("api_key")

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
	ch := chargehound.New("api_key")

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
	ch := chargehound.New("api_key")

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

		if b["template_id"] != "tmpl_1" {
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
		ID:         "dp_xxx",
		TemplateID: "tmpl_1",
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
	ch := chargehound.New("api_key")

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

		if b["template_id"] != "tmpl_1" {
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
		ID:         "dp_xxx",
		TemplateID: "tmpl_1",
		Products:   []chargehound.Product{{Name: "prod1", Amount: 100}},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDisputeCustomerInfo(t *testing.T) {
	ch := chargehound.New("api_key")

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

		if b["template_id"] != "tmpl_1" {
			t.Error("Incorrect template id.")
		}

		if b["customer_name"] != "name" {
			t.Error("Incorrect customer name.")
		}

		if b["customer_email"] != "email" {
			t.Error("Incorrect customer email.")
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
		ID:            "dp_xxx",
		TemplateID:    "tmpl_1",
		CustomerName:  "name",
		CustomerEmail: "email",
	})
	if err != nil {
		t.Error(err)
	}
}
