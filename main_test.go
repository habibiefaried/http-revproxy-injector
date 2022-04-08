package main

import (
	"bytes"
	"encoding/json"
	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func getWebHost() string {
	if os.Getenv("GITHUB_RUN_ID") != "" {
		return "http://revproxyinjector:4321"
	} else {
		return "http://localhost:4321"
	}
}

func getDVWAHost() string {
	if os.Getenv("GITHUB_RUN_ID") != "" {
		return "http://revproxydvwa:4322"
	} else {
		return "http://localhost:4322"
	}
}

func runCommand(t *testing.T, cmdin string) string {
	cmd := exec.Command("/bin/bash", "-c", cmdin)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(stdoutStderr))
		t.Fatal("Running command error: ", err)
	}

	return string(stdoutStderr)
}

func PostandCompare(t *testing.T, url string, values map[string]string, expected string) {
	jsonValue, err := json.Marshal(values)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
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

func DeleteandCompare(t *testing.T, url string, name string, expected string) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", url+"?name="+name, nil)
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
	PostandCompare(t, getWebHost()+"/revpr0xyconfig", values, "Data is injected")
	GetandCompare(t, getWebHost()+"/revpr0xyconfig", testData)
	GetandCompare(t, getWebHost()+"/revpr0xyconfig", "place\":\"header\"")
	DeleteandCompare(t, getWebHost()+"/revpr0xyconfig", cookieName, "Data is deleted")
}

func TestAccessProxy1(t *testing.T) {
	GetandCompare(t, getWebHost()+"/index.php", "Hello")
	GetandCompare(t, getWebHost()+"/info.php", "/usr/local/etc/php")
}

func TestAccessProxy2(t *testing.T) {
	GetandCompare(t, getWebHost()+"/headers.php", "X-Proxyinjector: In Action")
	testData := faker.Password()
	cookieName := faker.MonthName()
	values := map[string]string{"name": cookieName, "value": testData}
	PostandCompare(t, getWebHost()+"/revpr0xyconfig", values, "Data is injected")
	GetandCompare(t, getWebHost()+"/headers.php", cookieName+": "+testData)
}

func TestInsertConfig2(t *testing.T) {
	testData := faker.Password()
	values := map[string]string{"name": "", "value": testData}
	PostandCompare(t, getWebHost()+"/revpr0xyconfig", values, "variable name and value must exist")
}

func TestDeleteConfig(t *testing.T) {
	DeleteandCompare(t, getWebHost()+"/revpr0xyconfig", "randomized", "Data is deleted")
}

func TestInsertConfig3(t *testing.T) {
	testData := faker.Password()
	cookieName := faker.MonthName()
	values := map[string]string{"name": cookieName, "value": testData, "place": "randomized"}
	PostandCompare(t, getWebHost()+"/revpr0xyconfig", values, "variable 'place' is incorrect, should be 'header' or 'form' or 'query'")
}

func TestDVWA1(t *testing.T) {
	GetandCompare(t, getDVWAHost()+"/index.php", "Login :: Damn Vulnerable Web Application (DVWA)")
	values := map[string]string{"name": "Cookie", "value": "PHPSESSID=jv2db8n2jvjbjs4t44me934570; security=low", "place": "header"}
	PostandCompare(t, getDVWAHost()+"/revpr0xyconfig", values, "Data is injected")
	GetandCompare(t, getDVWAHost()+"/index.php", "vulnerabilities/sqli")
	DeleteandCompare(t, getDVWAHost()+"/revpr0xyconfig", "Cookie", "Data is deleted")
}

func TestSQLMapWithoutCookie(t *testing.T) {
	if os.Getenv("GITHUB_RUN_ID") != "" {
		output := runCommand(t, "sqlmap -u '"+getDVWAHost()+"/vulnerabilities/sqli/?id=2&Submit=Submit' -p id --dbs --batch")
		assert.Contains(t, output, "all tested parameters do not appear to be injectable")
	} else {
		t.Log("Skipped")
	}
}

func TestSQLMapWithCookie(t *testing.T) {
	if os.Getenv("GITHUB_RUN_ID") != "" {
		values := map[string]string{"name": "Cookie", "value": "PHPSESSID=jv2db8n2jvjbjs4t44me934570; security=low", "place": "header"}
		PostandCompare(t, getDVWAHost()+"/revpr0xyconfig", values, "Data is injected")
		output := runCommand(t, "sqlmap -u '"+getDVWAHost()+"/vulnerabilities/sqli/?id=2&Submit=Submit' -p id --dbs --batch")
		assert.Contains(t, output, "available databases [2]:")
	} else {
		t.Log("Skipped")
	}
}
