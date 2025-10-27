package services

import (
	"fmt"
	"net/http"

	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/repositories"
)

type HttpService struct {
	HTTP *repositories.HttpRepository
}

func NewHttpService(http *repositories.HttpRepository) *HttpService {
	return &HttpService{
		HTTP: http,
	}
}

func (s *HttpService) Get() (*entities.HttpResponse, error) {
	req, err := s.HTTP.GetHttpRequest(http.MethodGet, nil)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}
	httpResp, err := s.HTTP.GetHttpResponse(req)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}
	return httpResp, nil
}

func (s *HttpService) Update(method string, payloadObj interface{}) (*entities.HttpResponse, error) {
	payload, err := s.HTTP.GetHttpPayload(payloadObj)
	if err != nil {
		return nil, err
	}
	req, err := s.HTTP.GetHttpRequest(s.HTTP.RequestMethod(method), payload)
	if err != nil {
		return nil, err
	}
	httpResp, err := s.HTTP.GetHttpResponse(req)
	if err != nil {
		return nil, err
	}
	return httpResp, nil
}
