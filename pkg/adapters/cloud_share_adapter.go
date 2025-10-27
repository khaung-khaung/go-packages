package adapters

import (
	entities "github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/interfaces"
	"github.com/banyar/go-packages/pkg/repositories"
	"github.com/banyar/go-packages/pkg/services"
)

type CloudShareAdapter struct {
	CloudShareService interfaces.ICloudShareService
}

func NewCloudShareAdapter(DSNCloudShare *entities.DSNCloudShare) *CloudShareAdapter {
	cloudShareRepo := repositories.ConnectCloudShare(DSNCloudShare)
	return &CloudShareAdapter{
		CloudShareService: services.NewCloudShareService(cloudShareRepo),
	}
}
