package main

import "testing"

func TestParseKeyValue(t *testing.T) {
	kv := "blah    = asligfhsofddaodf"
	ParseKeyValue(kv)
}

func TestUnmarshal(t *testing.T) {
	input := `
[default]
aws_access_key_id = ASIAULREBYCTALZMWUMQ

aws_secret_access_key = Z6W26QvDY2se+PN6ooifEgbAZODpUlF0ZoLj+F5n

aws_session_token = IQoJb3JpZ2luX
`
	contents := Unmarshal([]byte(input))
	if contents == nil {
		t.Fail()
	}
}
