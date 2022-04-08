package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLogging(t *testing.T) {
	t.Log("Test success")
}

func TestInsertConfig1(t *testing.T) {
	values := map[string]string{"name": "Cookie", "value": "PHPSESSID=asdas$fa_aB-123134"}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post("http://revproxyinjector:4321/revpr0xyconfig", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("response Status: %v\n", resp.Status)
	t.Logf("response Headers: %v\n", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("response Body: %v\n", string(body))
}
