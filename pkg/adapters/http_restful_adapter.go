package adapters

import (
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"
)

type HttpRestfulAdapter struct {
	HttpService  interfaces.IHttpService
	BaseURL      string
	ExtraHeaders map[string]string
}

// NewHttpAdapter implementation
func NewHttpAdapter(baseURL string, extraHeaders map[string]string) *HttpRestfulAdapter {
	// Create HttpRepository with headers
	repo := repositories.ConnectRest(
		baseURL,
		repositories.WithExtraHeaders(extraHeaders),
	)

	return &HttpRestfulAdapter{
		HttpService:  services.NewHttpService(repo),
		BaseURL:      baseURL,
		ExtraHeaders: extraHeaders,
	}
}
