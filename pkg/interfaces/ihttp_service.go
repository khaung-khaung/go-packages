package interfaces

import (
	entities "github.com/banyar/go-packages/pkg/entities"
)

type IHttpService interface {
	Get() (*entities.HttpResponse, error)
	Update(method string, payloadObj interface{}) (*entities.HttpResponse, error)
}
