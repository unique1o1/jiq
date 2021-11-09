package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Failure = 2
	Success = 0
)

type mockRunner struct {
	returnSuccess bool
}

func (r mockRunner) run() int {
	if r.returnSuccess {
		return Success
	} else {
		return Failure
	}
}

// TODO: Test parsing args by breaking it out of main
// Until then, just run the mock
func TestJiqRun(t *testing.T) {
	var assert = assert.New(t)

	m := mockRunner{returnSuccess: true}
	result := doRun(m)
	assert.Equal(Success, result)
}

func TestJiqRunWithError(t *testing.T) {
	var assert = assert.New(t)

	m := mockRunner{returnSuccess: false}
	result := doRun(m)
	assert.Equal(Failure, result)
}
