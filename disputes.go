package chargehound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Wrapper for the Chargehound API disputes resource.
type Disputes struct {
	client *Client
}

// A dispute. See https://www.chargehound.com/docs/api/index.html#disputes.
type Dispute struct {
	ID                  string                 `json:"id"`
	State               string                 `json:"state"`
	Reason              string                 `json:"reason"`
	ChargedAt           string                 `json:"charged_at"`
	DisputedAt          string                 `json:"disputed_at"`
	DueBy               string                 `json:"due_by"`
	SubmittedAt         string                 `json:"submitted_at"`
	ClosedAt            string                 `json:"closed_at"`
	SubmittedCount      int                    `json:"submitted_count"`
	FileUrl             string                 `json:"file_url"`
	Template            string                 `json:"template"`
	Fields              map[string]interface{} `json:"fields"`
	MissingFields       map[string]interface{} `json:"missing_fields"`
	Products            []Product              `json:"products"`
	Charge              string                 `json:"charge"`
	IsChargeRefundable  bool                   `json:"is_charge_refundable"`
	Amount              int                    `json:"amount"`
	Currency            string                 `json:"currency"`
	Fee                 int                    `json:"fee"`
	ExternalCustomer    string                 `json:"external_customer"`
	CustomerName        string                 `json:"customer_name"`
	CustomerEmail       string                 `json:"customer_email"`
	CustomerPurchaseIp  string                 `json:"customer_purchase_ip"`
	AddressZip          string                 `json:"address_zip"`
	AddressLine1Check   string                 `json:"address_line1_check"`
	AddressZipCheck     string                 `json:"address_zip_check"`
	CvcCheck            string                 `json:"cvc_check"`
	StatementDescriptor string                 `json:"statement_descriptor"`
	Created             string                 `json:"created"`
	Updated             string                 `json:"updated"`
	Source              string                 `json:"source"`
}

// Dispute product data See https://www.chargehound.com/docs/api/index.html#product-data.
type Product struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	Sku         string `json:"sku,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	Url         string `json:"url,omitempty"`
}

// The type returned by a list disputes request. See https://www.chargehound.com/docs/api/index.html#retrieving-a-list-of-disputes.
type DisputeList struct {
	Data     []Dispute `json:"data"`
	HasMore  bool      `json:"has_more"`
	Livemode bool      `json:"livemode"`
	Object   string    `json:"object"`
	Url      string    `json:"url"`
}

// Params for a retrieve dispute request. See https://www.chargehound.com/docs/api/index.html#retrieving-a-dispute.
type RetrieveDisputeParams struct {
	// The dispute id.
	ID string
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client
}

// Params for a list disputes request. See https://www.chargehound.com/docs/api/index.html#retrieving-a-list-of-disputes.
type ListDisputesParams struct {
	Limit         int
	StartingAfter string
	EndingBefore  string
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client
}

// Params for updating or submitting a dispute. See https://www.chargehound.com/docs/api/index.html#updating-a-dispute.
type UpdateDisputeParams struct {
	// The dispute id.
	ID        string
	AccountID string
	Force     bool
	Template  string
	Charge    string
	Fields    map[string]interface{}
	Products  []Product
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client
}

type updateDisputeBody struct {
	Template  string                 `json:"template,omitempty"`
	Charge    string                 `json:"charge,omitempty"`
	AccountID string                 `json:"account_id,omitempty"`
	Force     bool                   `json:"force,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Products  []Product              `json:"products,omitempty"`
}

// Retrieve a single disputes.
func (dp *Disputes) Retrieve(params *RetrieveDisputeParams) (*Dispute, error) {
	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"GET",
		fmt.Sprintf("disputes/%s", params.ID),
		nil, // no body json
		nil, // no query params
	)

	if err != nil {
		return nil, err
	}

	var v Dispute
	err = req.newRequest(&v)

	return &v, err
}

// Retrieve a list of disputes.
func (dp *Disputes) List(params *ListDisputesParams) (*DisputeList, error) {
	// map the query params to a dict
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", string(params.Limit))
	}

	if params.StartingAfter != "" {
		q.Set("starting_after", params.StartingAfter)
	} else if params.EndingBefore != "" {
		q.Set("ending_before", params.EndingBefore)
	}

	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"GET",
		"disputes",
		nil, // no body json
		&q,
	)

	if err != nil {
		return nil, err
	}

	var v DisputeList
	err = req.newRequest(&v)

	return &v, err
}

func newUpdateDisputeBody(params *UpdateDisputeParams) (io.Reader, error) {
	body := updateDisputeBody{
		Fields:    params.Fields,
		Products:  params.Products,
		Template:  params.Template,
		AccountID: params.AccountID,
		Force:     params.Force,
		Charge:    params.Charge,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)

	return b, nil
}

// Update a dispute.
func (dp *Disputes) Update(params *UpdateDisputeParams) (*Dispute, error) {
	bodyJSON, err := newUpdateDisputeBody(params)
	if err != nil {
		return nil, err
	}

	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"PUT",
		fmt.Sprintf("disputes/%s", params.ID),
		bodyJSON,
		nil, // no query params
	)

	if err != nil {
		return nil, err
	}

	var v Dispute
	err = req.newRequest(&v)

	return &v, err
}

// Submit a dispute.
func (dp *Disputes) Submit(params *UpdateDisputeParams) (*Dispute, error) {
	bodyJSON, err := newUpdateDisputeBody(params)
	if err != nil {
		return nil, err
	}

	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"POST",
		fmt.Sprintf("disputes/%s/submit", params.ID),
		bodyJSON,
		nil, // no query params
	)

	if err != nil {
		return nil, err
	}

	var v Dispute
	err = req.newRequest(&v)

	return &v, err
}
