package main

import (
	"bytes"
	"encoding/json"
	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

var hostfortest string = "http://revproxyinjector:4321"

func PostandCompare(t *testing.T, values map[string]string, expected string) {
	jsonValue, err := json.Marshal(values)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(hostfortest+"/revpr0xyconfig", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), expected)
}

func GetandCompare(t *testing.T, fullurl string, expected string) {
	resp, err := http.Get(fullurl)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), expected)
}

func DeleteandCompare(t *testing.T, name string, expected string) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", hostfortest+"/revpr0xyconfig?name="+name, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Contains(t, string(respBody), expected)
}

func TestInsertConfig1(t *testing.T) {
	testData := faker.Password()
	cookieName := faker.MonthName()
	values := map[string]string{"name": cookieName, "value": testData}
	PostandCompare(t, values, "Data is injected")
	GetandCompare(t, hostfortest+"/revpr0xyconfig", testData)
	GetandCompare(t, hostfortest+"/revpr0xyconfig", "place\":\"header\"")
	DeleteandCompare(t, "Cookie", "Data is deleted")
}

func TestAccessProxy1(t *testing.T) {
	GetandCompare(t, hostfortest+"/index.php", "Hello")
	GetandCompare(t, hostfortest+"/info.php", "/usr/local/etc/php")
}

func TestAccessProxy2(t *testing.T) {
	GetandCompare(t, hostfortest+"/headers.php", "X-Proxyinjector: In Action")
	testData := faker.Password()
	cookieName := faker.MonthName()
	values := map[string]string{"name": cookieName, "value": testData}
	PostandCompare(t, values, "Data is injected")
	GetandCompare(t, hostfortest+"/headers.php", cookieName+": "+testData)
}

func TestInsertConfig2(t *testing.T) {
	testData := faker.Password()
	values := map[string]string{"name": "", "value": testData}
	PostandCompare(t, values, "variable name and value must exist")
}

func TestDeleteConfig(t *testing.T) {
	DeleteandCompare(t, "Cookie", "Data is deleted")
}

func TestInsertConfig3(t *testing.T) {
	testData := faker.Password()
	cookieName := faker.MonthName()
	values := map[string]string{"name": cookieName, "value": testData, "place": "randomized"}
	PostandCompare(t, values, "variable 'place' is incorrect, should be 'header' or 'form' or 'query'")
}
