package controller

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	OperationID string      `json:"operation_id"`
	Data        interface{} `json:"data"`
	MetaData    interface{} `json:"metadata"`
}

func WriteJSONResponse(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if response == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(response)
}
