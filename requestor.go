package chargehound

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type apiRequestor struct {
	apiKey      string
	userAgent   string
	bodyJSON    io.Reader
	httpClient  *http.Client
	method      string
	queryParams *url.Values
	url         string
}

func newAPIRequestor(cc *Client, optHTTP *http.Client, method, path string, bodyJSON io.Reader, queryParams *url.Values) (*apiRequestor, error) {
	var HTTPClient *http.Client

	if optHTTP != nil {
		HTTPClient = optHTTP
	} else {
		HTTPClient = cc.HTTPClient
	}

	url := cc.Protocol + cc.Host + cc.Basepath + path

	if queryParams != nil {
		url += "?" + queryParams.Encode()
	}

	requestor := apiRequestor{
		apiKey:      cc.ApiKey,
		bodyJSON:    bodyJSON,
		httpClient:  HTTPClient,
		method:      method,
		queryParams: queryParams,
		url:         url,
		userAgent:   "Chargehound/v1 GoBindings/" + cc.Version,
	}

	return &requestor, nil
}

func (ar *apiRequestor) newRequest(v interface{}) (*http.Response, error) {
	req, err := http.NewRequest(ar.method, ar.url, ar.bodyJSON)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(ar.apiKey, "")

	req.Header.Add("User-Agent", ar.userAgent)
	req.Header.Add("Content-Type", "application/json")

	res, err := ar.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, responseToError(res)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&v)
	return res, err
}
