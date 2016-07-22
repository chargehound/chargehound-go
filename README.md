# Chargehound Go bindings 

[![Build Status](https://travis-ci.org/chargehound/chargehound-go.svg?branch=master)](https://travis-ci.org/chargehound/chargehound-go)

## Installation

`go get -u github.com/chargehound/chargehound-go`

This library currently requires go >= 1.6.

## Usage

Every resource is accessed via the `Client` instance:

```go
ch := chargehound.New("{{your_api_key}}")

disputes, err := ch.Disputes.List()
```

## Documentation

[Disputes](https://www.chargehound.com/docs/api/index.html?go#disputes)

[Errors](https://test.chargehound.com/docs/api/index.html?go#errors)

## Development

Clone the latest source and run the tests:

```bash
$ git clone git@github.com:chargehound/chargehound-go.git
$ go test
```

Be sure to run gofmt on any code you plan on checking in.

## Deployment

To deploy a new version of the SDK, perform the following steps:

 1. Update the CHANGELOG to describe what features have been added.
 2. Bump the version number in `chargehound.go`.
 3. Create a tag for the version.
   ```git tag -a v{version} -m {message}```
 4. Push the tag to origin.
   ```git push origin v{version}```
 5. Create a release for the tag on [Github](https://help.github.com/articles/creating-releases/).
