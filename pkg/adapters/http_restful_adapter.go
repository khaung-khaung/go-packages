package adapters

import (
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"
)

type HttpRestfulAdapter struct {
	HttpService interfaces.IHttpService
	BaseURL     string
	Token       string
}

func NewHttpAdapter(baseURL string, token string) *HttpRestfulAdapter {
	http := repositories.ConnectRest(baseURL, token)
	return &HttpRestfulAdapter{
		HttpService: services.NewHttpService(http),
		BaseURL:     baseURL,
		Token:       token,
	}
}
