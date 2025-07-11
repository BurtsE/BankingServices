package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func insertUserIDToRequestBody(r *http.Request, uuid string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	newBody := map[string]interface{}{}
	err = json.Unmarshal(body, &newBody)
	if err != nil {
		return err
	}

	newBody["user_id"] = uuid

	modifiedBody, err := json.Marshal(newBody)
	if err != nil {
		return err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
	r.ContentLength = int64(len(modifiedBody))

	return nil
}
