package main

import (
	"testing"

	"go.uber.org/goleak"
	_ "go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
