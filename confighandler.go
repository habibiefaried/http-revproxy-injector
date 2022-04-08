package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		m.RLock()
		defer func() {
			m.RUnlock()
		}()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "OK", Data: &injecteddata})
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("%v", err)})
			return
		}
		var t RequestMessage
		err = json.Unmarshal(body, &t)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("%v", err)})
			return
		}

		name := t.Name
		value := t.Value
		place := t.Place

		if (name == "") || (value == "") {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "variable name and value must exist"})
			return
		}

		if place == "" {
			place = "header"
		} else {
			w.WriteHeader(402)
			if (place != "header") && (place != "form") && (place != "query") {
				json.NewEncoder(w).Encode(ResponseMessage{Message: "variable 'place' is incorrect, should be 'header' or 'form' or 'query'"})
				return
			}
		}

		m.Lock()
		injecteddata[name] = HTTPReqData{
			Value: value,
			Place: place,
		}
		m.Unlock()

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Data is injected"})
	case "DELETE":
		name, ok := r.URL.Query()["name"]

		if !ok || len(name[0]) < 1 {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "variable name must exist"})
			return
		}

		m.Lock()
		delete(injecteddata, name[0])
		m.Unlock()

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Data is deleted"})
	default:
		w.WriteHeader(504)
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
