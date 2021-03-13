package controller

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		response interface{}
	}{
		{
			"successful response",
			200,
			Response{
				OperationID: "foo_get",
				Data: map[string]string{
					"foo": "bar",
				},
				MetaData: map[string]string{
					"ping": "pong",
				},
			},
		},
		{
			"successful empty response",
			200,
			nil,
		},
		{
			"bad request",
			400,
			Response{
				OperationID: "foo_get",
				Data: map[string]string{
					"foo": "bar",
				},
				MetaData: map[string]string{
					"ping": "pong",
				},
			},
		},
		{
			"bad request empty response",
			400,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONResponse(w, tt.code, tt.response)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			assert.Equal(t, tt.code, resp.StatusCode)

			if tt.response != nil {
				expectedBody, _ := json.Marshal(tt.response)
				assert.JSONEq(t, string(expectedBody), string(body))
				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			} else {
				assert.Empty(t, string(body))
			}
		})
	}
}
