package controller

import (
	"net/http"
)

type WellKnownController struct {
	AppVersion    string
	AppCommitHash string
}

func (c *WellKnownController) Root(w http.ResponseWriter, _ *http.Request) {
	WriteJSONResponse(w, http.StatusOK, Response{
		OperationID: "app_get",
		Data: map[string]string{
			"service":     "kafka-test",
			"version":     c.AppVersion,
			"commit_hash": c.AppCommitHash,
		},
	})
}

func (c *WellKnownController) Health(w http.ResponseWriter, _ *http.Request) {
	WriteJSONResponse(w, http.StatusOK, Response{
		OperationID: "app_health",
		Data: map[string]string{
			"status": "OK",
		},
	})
}

func (c *WellKnownController) Version(w http.ResponseWriter, _ *http.Request) {
	WriteJSONResponse(w, http.StatusOK, map[string]string{
		"version":     c.AppVersion,
		"commit_hash": c.AppCommitHash,
	})
}
