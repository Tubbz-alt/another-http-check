package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestHTTPCodes(t *testing.T) {
	statusCodes := [4]int{200, 302, 404, 500}
	for _, statusCode := range statusCodes {
		r := &Request{
			Scheme:  "https",
			Host:    "httpbin.org",
			Port:    443,
			URI:     fmt.Sprintf("/status/%s", strconv.Itoa(statusCode)),
			Timeout: 30,
			Verbose: false,
		}

		var currrentStatusCodes []int
		currrentStatusCodes = append(currrentStatusCodes, statusCode)
		e := &Expected{
			StatusCodes: currrentStatusCodes,
		}

		msg, code, err := Check(r, e)

		if !strings.HasPrefix(msg, "OK") {
			t.Errorf("Wrong message [URI: %s]", r.URI)
		}

		if code != EXIT_OK {
			t.Errorf("Wrong exit code [URI: %s]", r.URI)
		}

		if err != nil {
			t.Errorf("Returned error is not nil [URI: %s]", r.URI)
		}
	}
}

func TestHTTPWrongCodes(t *testing.T) {
	statusCodes := [4]int{200, 302, 404, 500}
	for _, statusCode := range statusCodes {
		r := &Request{
			Scheme:  "https",
			Host:    "httpbin.org",
			Port:    443,
			URI:     fmt.Sprintf("/status/%s", strconv.Itoa(statusCode)),
			Timeout: 30,
			Verbose: false,
		}

		var currrentStatusCodes []int
		currrentStatusCodes = append(currrentStatusCodes, statusCode+1)
		e := &Expected{
			StatusCodes: currrentStatusCodes,
		}

		msg, code, err := Check(r, e)

		if !strings.HasPrefix(msg, "CRITICAL") {
			t.Errorf("Wrong message [URI: %s]", r.URI)
		}

		if code != EXIT_CRITICAL {
			t.Errorf("Wrong exit code [URI: %s]", r.URI)
		}

		if err != nil {
			t.Errorf("Returned error is not nil [URI: %s]", r.URI)
		}
	}
}

func TestTimeout(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/delay/10",
		Timeout: 5,
		Verbose: false,
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "CRITICAL") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_CRITICAL {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err == nil {
		t.Errorf("Returned error is nil [URI: %s]", r.URI)
	}

	if !strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
		t.Errorf("Non-timeout error renturned [URI: %s]", r.URI)
	}
}

func TestBasicAuthSuccess(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/basic-auth/user/password",
		Timeout: 30,
		Verbose: false,
		Authentication: Authentication{
			Type:     AUTH_BASIC,
			User:     "user",
			Password: "password",
		},
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "OK") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_OK {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err != nil {
		t.Errorf("Returned error is not nil [URI: %s]", r.URI)
	}
}

func TestBasicAuthFail(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/basic-auth/user/password",
		Timeout: 30,
		Verbose: false,
		Authentication: Authentication{
			Type:     AUTH_BASIC,
			User:     "user",
			Password: "password_",
		},
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "CRITICAL") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_CRITICAL {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err != nil {
		t.Errorf("Returned error is not nil [URI: %s]", r.URI)
	}
}

func TestContainsTextOK(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/anything?foobar=baz",
		Timeout: 30,
		Verbose: false,
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
		BodyText:    "foobar",
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "OK") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_OK {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err != nil {
		t.Errorf("Returned error is not nil [URI: %s]", r.URI)
	}
}

func TestContainsTextFail(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/anything?foobar=baz",
		Timeout: 30,
		Verbose: false,
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
		BodyText:    "loremipsum",
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "CRITICAL") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_CRITICAL {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err != nil {
		t.Errorf("Returned error is not nil [URI: %s]", r.URI)
	}
}

func TestSSLOK(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/anything",
		Timeout: 30,
		Verbose: false,
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
		SSLCheck: SSLCheck{
			Run:          true,
			DaysWarning:  20,
			DaysCritical: 5,
		},
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "OK") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_OK {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err != nil {
		t.Errorf("Returned error is not nil [URI: %s]", r.URI)
	}
}

func TestSSLWarning(t *testing.T) {
	r := &Request{
		Scheme:  "https",
		Host:    "httpbin.org",
		Port:    443,
		URI:     "/anything",
		Timeout: 30,
		Verbose: false,
	}

	var currrentStatusCodes []int
	currrentStatusCodes = append(currrentStatusCodes, 200)
	e := &Expected{
		StatusCodes: currrentStatusCodes,
		SSLCheck: SSLCheck{
			Run:          true,
			DaysWarning:  10000,
			DaysCritical: 5,
		},
	}

	msg, code, err := Check(r, e)

	if !strings.HasPrefix(msg, "WARNING") {
		t.Errorf("Wrong message [URI: %s]", r.URI)
	}

	if code != EXIT_WARNING {
		t.Errorf("Wrong exit code [URI: %s]", r.URI)
	}

	if err != nil {
		t.Errorf("Returned error is not nil [URI: %s]", r.URI)
	}
}