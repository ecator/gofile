package server

import (
	"encoding/json"
	"net/http"
)

type jsonData struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}

func resJson(w http.ResponseWriter, jData *jsonData) {
	b, _ := json.Marshal(jData)
	w.Write(b)
}
