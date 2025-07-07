package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func insertIDToRequestBody(r *http.Request, uuid string) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	newBody := map[string]interface{}{}
	err = json.Unmarshal(body, &newBody)
	if err != nil {
		return err
	}

	newBody["uuid"] = uuid

	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}
