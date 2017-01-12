// Testing for http handlers found in common.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"log"
)

func testStatus(expected int, rr *httptest.ResponseRecorder, t *testing.T) {
	if expected != rr.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, expected)
	}
}

func testBody(expected string, rr *httptest.ResponseRecorder, t *testing.T) {
	if expected != rr.Body.String() {
		t.Errorf("handler returned unexpected body: got %s want %v", rr.Body.String(), expected)
	}
}

func setup() (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", "/not_found_url", nil)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()

	return req, rr

}

func TestNotFoundHandler0(t *testing.T) {
	req, rr := setup()

	NotFoundHandler(rr, req, "")

	testStatus(StatusNotFound, rr, t)

	expected := "{ \"error_message\": \"Failure to retrieve objects\"}"
	testBody(expected, rr, t)

}

func TestNotFoundHandler1(t *testing.T) {
	req, rr := setup()

	errMsg := "2Do not found."
	NotFoundHandler(rr, req, errMsg)

	testStatus(StatusNotFound, rr, t)

	expected := "{\"error_message\":\"2Do not found.\"}"
	testBody(expected, rr, t)
}

func TestInternalErrorHandler0(t *testing.T) {
	req, rr := setup()

	InternalErrorHandler(rr, req, "")

	testStatus(StatusInternalError, rr, t)

	expected := "{ \"error_message\": \"The server encountered an unexpected condition which prevented it from fulfilling the request.\"}"
	testBody(expected, rr, t)

}

func TestInternalErrorHandler1(t *testing.T) {
	req, rr := setup()

	errMsg := "Internal Error"
	InternalErrorHandler(rr, req, errMsg)

	testStatus(StatusInternalError, rr, t)

	expected := fmt.Sprintf("{\"error_message\":\"%s\"}", errMsg)
	testBody(expected, rr, t)
}

func TestUnauthorizedHandler0(t *testing.T) {
	req, rr := setup()

	UnauthorizedHandler(rr, req, "")

	testStatus(StatusUnauthorized, rr, t)

	expected := "{ \"error_message\": \"The request requires user authentication.\"}"
	testBody(expected, rr, t)
}

func TestUnauthorizedHandler1(t *testing.T) {
	req, rr := setup()

	errMsg := "Unauthorized error"
	UnauthorizedHandler(rr, req, errMsg)

	testStatus(StatusUnauthorized, rr, t)

	expected := fmt.Sprintf("{\"error_message\":\"%s\"}", errMsg)
	testBody(expected, rr, t)
}

func TestBadRequestHandler0(t *testing.T) {
	req, rr := setup()

	BadRequestHandler(rr, req, "")

	testStatus(StatusBadRequest, rr, t)

	expected := "{ \"error_message\": \"Bad Request.\"}"
	testBody(expected, rr, t)
}

func TestBadRequestHandler1(t *testing.T) {
	req, rr := setup()

	errMsg := "Bad Request response"
	UnauthorizedHandler(rr, req, errMsg)

	testStatus(StatusUnauthorized, rr, t)

	expected := fmt.Sprintf("{\"error_message\":\"%s\"}", errMsg)
	testBody(expected, rr, t)
}
