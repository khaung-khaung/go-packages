package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.frontiir.net/sa-dev/log-package/pkg/frontlog"
	entities "github.com/banyar/go-packages/pkg/entities"
	"go.uber.org/zap"
)

type HttpRepository struct {
	baseURL      string
	extraHeaders map[string]string
}

type ConfigOption func(*HttpRepository)

func ConnectRest(baseURL string, opts ...ConfigOption) *HttpRepository {
	repo := &HttpRepository{
		baseURL: baseURL,
	}
	// Apply all provided options
	for _, opt := range opts {
		opt(repo)
	}
	return repo
}

func WithExtraHeaders(headers map[string]string) ConfigOption {
	return func(h *HttpRepository) {
		if h.extraHeaders == nil {
			h.extraHeaders = make(map[string]string)
		}
		for k, v := range headers {
			h.extraHeaders[k] = v
		}
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
		// Timeout: 10000, // Use configured timeout
	}
	// Set the Authorization header
	for k, v := range s.extraHeaders {
		req.Header.Add(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	// req.Header.Add("Authorization", s.token)
	// fmt.Println("http req =====================> ", req)
	// Make the request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}

	defer resp.Body.Close() //

	s.logResponse(resp)

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

func (s *HttpRepository) logResponse(resp *http.Response) {
	// Read body preview (max 1KB)
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	resp.Body = io.NopCloser(io.MultiReader(
		bytes.NewReader(bodyBytes),
		resp.Body,
	))

	// Prepare headers map
	headers := make(map[string]string)
	for key, values := range resp.Header {
		for _, value := range values {
			if strings.EqualFold(key, "Authorization") {
				value = "[REDACTED]"
			}
			headers[key] = value
		}
	}

	// Log with zap
	frontlog.Logger.Info("[HTTP RESPONSE]",
		zap.String("status", resp.Status),
		zap.Int("status_code", resp.StatusCode),
		zap.Any("headers", headers),
		// zap.String("body_preview", string(bodyBytes)),
	)
}
