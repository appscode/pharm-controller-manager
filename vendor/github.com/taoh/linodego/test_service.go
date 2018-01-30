package linodego

import (
	"encoding/json"
	"net/url"
)

// TestService is an interface generated for "github.com/appscode/linodego".TestService.
type TestInterface interface {
	Echo(string, string, *TestResponse) error
}

// Test Service
type TestService struct {
	client *Client
}

var _ TestInterface = &TestService{}

// Test Service Response
type TestResponse struct {
	Response
	Data map[string]string
}

// Echo request with the given key and value
func (t *TestService) Echo(key string, val string, v *TestResponse) error {
	u := &url.Values{}
	u.Add(key, val)
	if err := t.client.do("test.echo", u, &v.Response); err != nil {
		return err
	}
	v.Data = map[string]string{}
	if err := json.Unmarshal(v.RawData, &v.Data); err != nil {
		return err
	}
	return nil
}
