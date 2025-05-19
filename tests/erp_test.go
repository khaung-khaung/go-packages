package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"
)

type ERPPlanDetail struct {
	LocalServiceId     string `json:"local_service_id"` // Maps directly to JSON "id"
	PlanStartDate      string `json:"plan_start_date"`
	SubscribedPlan     string `json:"subscribed_plan"`
	SubscriptionStatus string `json:"subscription_status"`
}

func (d *ERPPlanDetail) UnmarshalJSON(data []byte) error {
	// Create an alias type to avoid recursion
	type Alias ERPPlanDetail

	// Define auxiliary struct that embeds the alias
	aux := &struct {
		ID    string `json:"id"`
		Dates struct {
			StartDate string `json:"start_date"`
		} `json:"dates"`
		ServiceDetails struct {
			Package string `json:"package"`
		} `json:"service_details"`
		Status struct {
			State string `json:"state"`
		} `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Map nested fields to flat structure
	d.PlanStartDate = aux.Dates.StartDate
	d.SubscribedPlan = aux.ServiceDetails.Package
	d.SubscriptionStatus = aux.Status.State
	d.LocalServiceId = aux.ID

	return nil
}

func TestPlanSummary(t *testing.T) {
	cid := "SCP000012"
	rangeOption := "[0,10]"

	baseURL := os.Getenv("ERP_BASE_URL")
	apiToken := os.Getenv("ERP_TOKEN")
	apiUrl := fmt.Sprintf("%s/subscriptions?range=%s&filter={\"device_id\":\"%s\"}", baseURL, rangeOption, cid)

	httpAdapter := adapters.NewHttpAdapter(apiUrl, apiToken)
	fmt.Println("apiUrl", httpAdapter.BaseURL)
	fmt.Println("apiToken", httpAdapter.Token)

	resp, err := httpAdapter.HttpService.Get()
	if err != nil {
		log.Fatal("ERROR : ", err)
	}
	common.DisplayJsonFormat("rest2Adapter", resp)
}

func TestPlanDetail(t *testing.T) {
	serviceId := "2061125-001"

	baseURL := os.Getenv("ERP_BASE_URL")
	apiToken := os.Getenv("ERP_TOKEN")
	apiUrl := fmt.Sprintf("%s/subscriptions/%s", baseURL, serviceId)
	httpAdapter := adapters.NewHttpAdapter(apiUrl, apiToken)
	fmt.Println("apiUrl", httpAdapter.BaseURL)
	fmt.Println("apiToken", httpAdapter.Token)

	resp, err := httpAdapter.HttpService.Get()
	if err != nil {
		log.Fatal("ERROR : ", err)
	}
	var planDetail *ERPPlanDetail
	if err := json.Unmarshal([]byte(resp.Message), &planDetail); err != nil {
		log.Println("Error parsing JSON:", err)

	}

	common.DisplayJsonFormat("rest2Adapter", planDetail)
}
