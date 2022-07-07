package main

import "testing"

func TestParseKeyValue(t *testing.T) {
	kv := "blah    = asligfhsofddaodf"
	ParseKeyValue(kv)
}
