package yopass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPClient allows modifying the underlying http.Client.
var HTTPClient = http.DefaultClient

// ServerError represents a yopass server error.
type ServerError struct {
	err error
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("yopass server error: %s", e.err)
}

func (e *ServerError) Unwrap() error {
	return e.err
}

type serverResponse struct {
	Message string `json:"message"`
	OneTime bool   `json:"one_time"`
}

// Fetch retrieves a secret by its ID from the specified server.
func Fetch(server string, id string) (string, bool, error) {
	server = strings.TrimSuffix(server, "/")

	resp, err := HTTPClient.Get(server + "/secret/" + id)
	if err != nil {
		return "", false, &ServerError{err: err}
	}

	return handleServerResponse(resp)
}

func Delete(server string, id string) (string, error) {
	server = strings.TrimSuffix(server, "/")

	req, err := http.NewRequest("DELETE", server+"/secret/"+id, nil)
	if err != nil {
		panic(err)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "", &ServerError{err: err}
	}

	msg, _, err := handleServerResponse(resp)
	return msg, err
}

// Store sends the secret to the specified server and returns the secret ID.
func Store(server string, s Secret) (string, error) {
	server = strings.TrimSuffix(server, "/")

	var j bytes.Buffer
	if err := (json.NewEncoder(&j)).Encode(&s); err != nil {
		return "", fmt.Errorf("could not encode request: %w", err)
	}
	resp, err := HTTPClient.Post(server+"/secret", "application/json", &j)
	if err != nil {
		return "", &ServerError{err: err}
	}

	msg, _, err := handleServerResponse(resp)
	return msg, err
}

func handleServerResponse(resp *http.Response) (string, bool, error) {
	defer resp.Body.Close()

	var r serverResponse

	// Delete Operation
	if (resp.Request.Method == http.MethodDelete) && (resp.StatusCode == http.StatusNoContent) { // During DELETE Call, status code will be 204 (successful deletion)
		// As response body will be empty after delete, upon successful deletion response will be empty
		return "", false, nil
	}

	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(msg, &r); err == nil {
			msg = []byte(r.Message)
		}
		err := fmt.Errorf("%s: unexpected response %s: %s", resp.Request.Method, resp.Status, string(msg))
		return "", false, &ServerError{err: err}
	}


	if err := (json.NewDecoder(resp.Body)).Decode(&r); err != nil {
		return "", false, fmt.Errorf("could not decode server response: %w", err)
	}


	return r.Message, r.OneTime, nil
}
