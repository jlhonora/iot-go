package main

import (
	"encoding/json"
	"testing"
)

func TestNodeDecode(t *testing.T) {
	in := []byte(`{"id":1,"identifier":"1234","name":"TestName","created_at":"2014-10-20 10:01:04","updated_at": "2014-10-20 10:01:04"}`)
	var out Node
	err := json.Unmarshal(in, &out)
	if err != nil {
		t.Error(err)
	}
	if out.Id != 1 {
		t.Error("Expected Id 1, got ", out.Id)
	}
	if out.Name != "TestName" {
		t.Error("Expected Name TestName, got ", out.Name)
	}
}
