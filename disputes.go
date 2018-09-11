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
	// A unique identifier for the dispute. This id is set by the payment processor of the dispute.
	ID string `json:"id"`
	// State of the dispute. One of `needs_response`,`submitted`, `under_review`, `won`, `lost`, `warning_needs_response`, `warning_under_review`, `warning_closed` , `response_disabled`, `charge_refunded`, `accepted`, `queued`.
	State string `json:"state"`
	// Reason for the dispute. One of `fraudulent`, `unrecognized`, `general`, `duplicate`, `subscription_canceled`, `product_unacceptable`, `product_not_received`, `credit_not_processed`, `incorrect_account_details`, `insufficient_funds`, `bank_cannot_process`, `debit_not_authorized`, `goods_services_returned_or_refused`, `goods_services_cancelled` |
	Reason string `json:"reason"`
	// ISO 8601 timestamp - when the charge was made.
	ChargedAt string `json:"charged_at"`
	// ISO 8601 timestamp - when the charge was disputed.
	DisputedAt string `json:"disputed_at"`
	// ISO 8601 timestamp - when dispute evidence needs to be disputed by.
	DueBy string `json:"due_by"`
	// ISO 8601 timestamp - when dispute evidence was submitted.
	SubmittedAt string `json:"submitted_at"`
	// ISO 8601 timestamp - when the dispute was resolved.
	ClosedAt string `json:"closed_at"`
	// Number of times the dispute evidence has been submitted.
	SubmittedCount int `json:"submitted_count"`
	// Id of the template attached to the dispute.
	Template string `json:"template"`
	// Evidence fields attached to the dispute.
	Fields map[string]interface{} `json:"fields"`
	// Any fields required by the template that have not yet been provided.
	MissingFields map[string]interface{} `json:"missing_fields"`
	// A list of products in the disputed order. (See [Product data](#product-data) for details.) (optional)
	Products []Product `json:"products"`
	// Id of the disputed charge.
	Charge string `json:"charge"`
	// Can the charge be refunded.
	IsChargeRefundable bool `json:"is_charge_refundable"`
	// Amount of the disputed charge. Amounts are in cents (or other minor currency unit.)
	Amount int `json:"amount"`
	// Currency code of the disputed charge. e.g. 'USD'.
	Currency string `json:"currency"`
	// Dispute fee.
	Fee int `json:"fee"`
	// The amount deducted due to the chargeback. Amounts are in cents (or other minor currency unit.)
	ReversalAmount int `json:"reversal_amount"`
	// Currency code of the deduction amount. e.g. 'USD'.
	ReversalCurrency string `json:"reversal_currency"`
	// Id of the customer (if any). This id is set by the payment processor of the dispute.
	Customer string `json:"customer"`
	// Name of the customer (if any).
	CustomerName string `json:"customer_name"`
	// Email of the customer (if any).
	CustomerEmail string `json:"customer_email"`
	// IP of purchase (if available).
	CustomerPurchaseIP string `json:"customer_purchase_ip"`
	// Billing address zip of the charge.
	AddressZip string `json:"address_zip"`
	// State of address check (if available). One of `pass`, `fail`, `unavailable`, `checked`.
	AddressLine1Check string `json:"address_line1_check"`
	// State of address zip check (if available). One of `pass`, `fail`, `unavailable`, `checked`.
	AddressZipCheck string `json:"address_zip_check"`
	// State of cvc check (if available). One of `pass`, `fail`, `unavailable`, `checked`.
	CVCCheck string `json:"cvc_check"`
	// The descriptor that appears on the customer's credit card statement for this change.
	StatementDescriptor string `json:"statement_descriptor"`
	// The account id for Connected accounts that are charged directly through Stripe (if any)
	UserID string `json:"user_id"`
	// The kind for the dispute, 'chargeback', 'retrieval' or 'pre_arbitration'.
	Kind string `json:"kind"`
	// ISO 8601 timestamp.
	Created string `json:"created"`
	// ISO 8601 timestamp.
	Updated string `json:"updated"`
	// The source of the dispute. One of `mock`, `braintree`, `api` or `stripe`
	Source string `json:"source"`
	// The payment processor of the dispute. One of `braintree` or `stripe`
	Processor string `json:"processor"`
	// Custom URL with dispute information.
	ReferenceURL string `json:"reference_url"`
	// Data about the API response that created dispute.
	Response HTTPResponse `json:"-"`
}

// Dispute product data See https://www.chargehound.com/docs/api/index.html#product-data.
type Product struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	Sku         string `json:"sku,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	URL         string `json:"url,omitempty"`
}

// The type returned by a list disputes request. See https://www.chargehound.com/docs/api/index.html#retrieving-a-list-of-disputes.
type DisputeList struct {
	Data     []Dispute    `json:"data"`
	HasMore  bool         `json:"has_more"`
	Livemode bool         `json:"livemode"`
	Object   string       `json:"object"`
	URL      string       `json:"url"`
	Response HTTPResponse `json:"-"`
}

// The type returned by a dispute response request.
type Response struct {
	Livemode       bool                   `json:"livemode"`
	DisputeID      string                 `json:"dispute_id"`
	ExternalCharge string                 `json:"external_charge"`
	AccountID      string                 `json:"account_id"`
	Evidence       map[string]interface{} `json:"evidence"`
	ResponseURL    string                 `json:"response_url"`
	Response       HTTPResponse           `json:"-"`
}

// Params for a retrieve dispute request. See https://www.chargehound.com/docs/api/index.html#retrieving-a-dispute.
type RetrieveDisputeParams struct {
	// The dispute id.
	ID string
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client
}

// Params for a dispute accept request.
type AcceptDisputeParams struct {
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
	State         string
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client
}

// Data about the API response that created the Chargehound object.
type HTTPResponse struct {
	// The HTTP status code.
	Status int
}

// Params for updating or submitting a dispute. See https://www.chargehound.com/docs/api/index.html#updating-a-dispute.
type UpdateDisputeParams struct {
	// The dispute id.
	ID     string
	UserID string
	// Id of the connected account for this dispute (if multiple accounts are connected)
	Account      string
	Force        bool
	Queue        bool
	Template     string
	Charge       string
	Fields       map[string]interface{}
	Products     []Product
	ReferenceURL string
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client
}

type CreateDisputeParams struct {
	// The id of the dispute in your payment processor. For Stripe looks like `dp_XXX`.
	ID string `json:"id"`
	// The id of the disputed charge in your payment processor. For Stripe looks like `ch_XXX`.
	Charge string `json:"charge"`
	// The id of the charged customer in your payment processor. For Stripe looks like `cus_XXX`. (optional)
	Customer string `json:"customer,omitempty"`
	// The bank provided reason for the dispute. One of `general`, `fraudulent`, `duplicate`, `subscription_canceled`, `product_unacceptable`, `product_not_received`, `unrecognized`, `credit_not_processed`, `incorrect_account_details`, `insufficient_funds`, `bank_cannot_process`, `debit_not_authorized`.
	Reason string `json:"reason"`
	// ISO 8601 timestamp - when the charge was made.
	ChargedAt string `json:"charged_at"`
	// ISO 8601 timestamp - when the charge was disputed.
	DisputedAt string `json:"disputed_at"`
	// ISO 8601 timestamp - when dispute evidence needs to be disputed by.
	DueBy string `json:"due_by"`
	// The currency code of the disputed charge. e.g. 'USD'.
	Currency string `json:"currency"`
	// The amount of the disputed charge. Amounts are in cents (or other minor currency unit.)
	Amount int `json:"amount"`
	// The payment processor for the charge. One of `braintree` or `stripe`. (optional)
	Processor string `json:"processor,omitempty"`
	// The state of the dispute. One of `needs_response`, `warning_needs_response`. (optional)
	State string `json:"state,omitempty"`
	// The currency code of the dispute balance withdrawal. e.g. 'USD'. (optional)
	ReversalCurrency string `json:"reversal_currency,omitempty"`
	// The amount of the dispute fee. Amounts are in cents (or other minor currency unit.) (optional)
	Fee int `json:"fee,omitempty"`
	// The amount of the dispute balance withdrawal (without fee). Amounts are in cents (or other minor currency unit.) (optional)
	ReversalAmount int `json:"reversal_amount,omitempty"`
	// The total amount of the dispute balance withdrawal (with fee). Amounts are in cents (or other minor currency unit.) (optional)
	ReversalTotal int `json:"reversal_total,omitempty"`
	// Is the disputed charge refundable. (optional)
	IsChargeRefundable bool `json:"is_charge_refundable,omitempty"`
	// How many times has dispute evidence been submitted. (optional)
	SubmittedCount int `json:"submitted_count,omitempty"`
	// State of address check (if available). One of `pass`, `fail`, `unavailable`, `checked`. (optional)
	AddressLine1Check string `json:"address_line1_check,omitempty"`
	// State of address zip check (if available). One of `pass`, `fail`, `unavailable`, `checked`. (optional)
	AddressZipCheck string `json:"address_zip_check,omitempty"`
	// State of cvc check (if available). One of `pass`, `fail`, `unavailable`, `checked`. (optional)
	CVCCheck string `json:"cvc_check,omitempty"`
	// The id of the template to use. (optional)
	Template string `json:"template,omitempty"`
	// Key value pairs to hydrate the template's evidence fields. (optional)
	Fields map[string]interface{} `json:"fields,omitempty"`
	// List of products the customer purchased. (optional)
	Products []Product `json:"products,omitempty"`
	// Set the account id for Connected accounts that are charged directly through Stripe. (optional)
	UserID string `json:"user_id,omitempty"`
	// Set the kind for the dispute, 'chargeback', 'retrieval' or 'pre_arbitration'. (optional)
	Kind string `json:"kind,omitempty"`
	// Submit dispute evidence immediately after creation. (optional)
	Submit bool `json:"submit,omitempty"`
	// Queue dispute for submission immediately after creation. (optional)
	Queue bool `json:"queue,omitempty"`
	// Custom URL with dispute information.
	ReferenceURL string `json:"reference_url,omitempty"`
	// Optional http client for the request. Typically needed when using App Engine.
	OptHTTPClient *http.Client `json:"-"`
}

type updateDisputeBody struct {
	Template     string                 `json:"template,omitempty"`
	Charge       string                 `json:"charge,omitempty"`
	Account      string                 `json:"account,omitempty"`
	UserID       string                 `json:"user_id,omitempty"`
	ReferenceURL string                 `json:"reference_url,omitempty"`
	Force        bool                   `json:"force,omitempty"`
	Queue        bool                   `json:"queue,omitempty"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Products     []Product              `json:"products,omitempty"`
}

// Create a dispute
func (dp *Disputes) Create(params *CreateDisputeParams) (*Dispute, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(params)

	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"POST",
		"disputes",
		b,
		nil, // no query params
	)

	if err != nil {
		return nil, err
	}

	var v Dispute
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

	return &v, err
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
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

	return &v, err
}

// Retrieve the response for a dispute.
func (dp *Disputes) Response(params *RetrieveDisputeParams) (*Response, error) {
	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"GET",
		fmt.Sprintf("disputes/%s/response", params.ID),
		nil, // no body json
		nil, // no query params
	)

	if err != nil {
		return nil, err
	}

	var v Response
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

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

	if params.State != "" {
		q.Set("state", params.State)
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
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

	return &v, err
}

func newUpdateDisputeBody(params *UpdateDisputeParams) (io.Reader, error) {
	body := updateDisputeBody{
		Fields:       params.Fields,
		Products:     params.Products,
		ReferenceURL: params.ReferenceURL,
		Template:     params.Template,
		UserID:       params.UserID,
		Account:      params.Account,
		Force:        params.Force,
		Queue:        params.Queue,
		Charge:       params.Charge,
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
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

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
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

	return &v, err
}

// Accept a dispute.
func (dp *Disputes) Accept(params *AcceptDisputeParams) (*Dispute, error) {
	req, err := newAPIRequestor(
		dp.client,
		params.OptHTTPClient,
		"POST",
		fmt.Sprintf("disputes/%s/accept", params.ID),
		nil, // no body json
		nil, // no query params
	)

	if err != nil {
		return nil, err
	}

	var v Dispute
	res, err := req.newRequest(&v)
	if err == nil {
		v.Response = HTTPResponse{Status: res.StatusCode}
	}

	return &v, err
}
