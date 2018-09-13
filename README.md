# Chargehound Go bindings 
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/chargehound/chargehound-go) [![Build Status](https://travis-ci.org/chargehound/chargehound-go.svg?branch=master)](https://travis-ci.org/chargehound/chargehound-go)

## Installation

`go get github.com/chargehound/chargehound-go`

This library currently requires go >= 1.6.

## Usage

Every resource is accessed via the `Client` instance:

```go
ch := chargehound.New("{{your_api_key}}", nil)
```

### Requests

Go requests use defined structs to represent parameters.

```go
params := chargehound.UpdateDisputeParams{
  ID:       "dp_123",
  Template: "unrecognized",
  Fields: map[string]interface{}{
    "customer_name": "Susie Chargeback",
  },
}

dispute, err := ch.Disputes.Submit(&params)
```

### Responses

Responses from the API are automatically parsed from JSON and returned as Go structs.

Responses also include the HTTP status code on the response object as the status field.

```go
dispute, err := ch.Disputes.Retrieve("dp_123")

fmt.Println("{}", dispute.State)
// "needs_response"
fmt.Println("{}", dispute.Response.Status)
// 200
```

## Documentation

[Disputes](https://www.chargehound.com/docs/api/index.html?go#disputes)

[Errors](https://www.chargehound.com/docs/api/index.html?go#errors)

The Go library returns adapted structs rather than JSON from API calls.

## Google AppEngine

If you're using the library in a Google App Engine environment, you can pass a custom http client along with each request. `OptHTTPClient` is defined on all param structs.

```go
import (
  "net/http"

  "google.golang.org/appengine"
  "google.golang.org/appengine/urlfetch"

  "github.com/chargehound/chargehound-go"
)

ch := chargehound.New("{{your_api_key}}", nil)

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  httpClient := urlfetch.Client(c)

  params := chargehound.ListDisputeParams{
    OptHTTPClient: &httpClient
  }

  disputes, err := ch.Disputes.List(&params)
}
```

## Development

Clone the latest source and run the tests:

```bash
$ git clone git@github.com:chargehound/chargehound-go.git
$ go test
```

Be sure to run gofmt on any code you plan on checking in.

```bash
gofmt -s -w .
```

## Deployment

To deploy a new version of the SDK, perform the following steps:

 1. Update the CHANGELOG to describe what features have been added.
 2. Bump the version number in `chargehound.go`.
 3. Create a tag for the version.
   ```git tag -a v{version} -m {message}```
 4. Push the tag to origin.
   ```git push origin v{version}```
