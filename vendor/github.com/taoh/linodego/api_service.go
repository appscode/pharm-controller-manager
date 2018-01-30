package linodego

import (
	"encoding/json"
	"net/url"
)

// ApiService is an interface generated for "github.com/appscode/linodego".ApiService.
type ApiInterface interface {
	Spec(*ApiResponse) error
}

// API Service
type ApiService struct {
	client *Client
}

var _ ApiInterface = &ApiService{}

// Response for api.spec Service
type ApiResponse struct {
	Response
	Data map[string]interface{}
}

// Get API Specs
func (t *ApiService) Spec(v *ApiResponse) error {
	u := &url.Values{}
	if err := t.client.do("api.spec", u, &v.Response); err != nil {
		return err
	}
	v.Data = map[string]interface{}{}
	if err := json.Unmarshal(v.RawData, &v.Data); err != nil {
		return err
	}
	return nil
}
