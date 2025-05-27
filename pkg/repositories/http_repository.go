package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	entities "github.com/banyar/go-packages/pkg/entities"
)

type HttpRepository struct {
	baseURL string
	token   string
}

func ConnectRest(baseURL string, token string) *HttpRepository {
	return &HttpRepository{
		baseURL: baseURL,
		token:   token,
	}
}

func (s *HttpRepository) GetHttpPayload(payloadObj interface{}) (*bytes.Buffer, error) {
	payload, err := json.Marshal(payloadObj)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}
	return bytes.NewBuffer(payload), err
}

func (s *HttpRepository) GetHttpRequest(method string, payloadObj io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, s.baseURL, payloadObj)
	fmt.Println("HTTP REQUEST", req)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}
	return req, err
}

func (s *HttpRepository) RequestMethod(method string) string {
	var requestMethod = map[string]string{
		"PUT":  http.MethodPut,
		"POST": http.MethodPost,
	}
	return requestMethod[strings.ToUpper(method)]
}

func (s *HttpRepository) GetHttpResponse(req *http.Request) (*entities.HttpResponse, error) {
	// Create a new http.Client object
	client := &http.Client{
		Timeout: 20 * time.Second, // Set timeout here
	}

	// Set the Authorization header
	req.Header.Add("Authorization", s.token)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}

	var httpResponse *entities.HttpResponse
	// Get the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}

	httpResponse = &entities.HttpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
	defer resp.Body.Close()
	return httpResponse, nil
}
