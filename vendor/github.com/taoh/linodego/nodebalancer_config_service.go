package linodego

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// NodeBalancer Config Service
type NodeBalancerConfigService struct {
	client *Client
}

// Response for nodebalancer.config.list API
type NodeBalancerConfigListResponse struct {
	Response
	NodeBalancerConfigs []NodeBalancerConfig
}

// Response for general config APIs
type NodeBalancerConfigResponse struct {
	Response
	NodeBalancerConfigId NodeBalancerConfigId
}

// Get Config List. If configId is greater than 0, limit results to given config.
func (t *NodeBalancerConfigService) List(nodeBalancerId int, configId int) (*NodeBalancerConfigListResponse, error) {
	u := &url.Values{}
	if configId > 0 {
		u.Add("ConfigID", strconv.Itoa(configId))
	}
	u.Add("NodeBalancerID", strconv.Itoa(nodeBalancerId))
	v := NodeBalancerConfigListResponse{}
	if err := t.client.do("nodebalancer.config.list", u, &v.Response); err != nil {
		return nil, err
	}

	v.NodeBalancerConfigs = make([]NodeBalancerConfig, 5)
	if err := json.Unmarshal(v.RawData, &v.NodeBalancerConfigs); err != nil {
		return nil, err
	}
	return &v, nil
}

// Create Config
func (t *NodeBalancerConfigService) Create(nodeBalancerId int, args map[string]string) (*NodeBalancerConfigResponse, error) {
	u := &url.Values{}
	u.Add("NodeBalancerID", strconv.Itoa(nodeBalancerId))
	// add optional parameters
	processOptionalArgs(args, u)
	v := NodeBalancerConfigResponse{}
	if err := t.client.do("nodebalancer.config.create", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerConfigId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update Config. See https://www.linode.com/api/nodebalancer/nodebalancer.config.update for allowed arguments.
func (t *NodeBalancerConfigService) Update(configId int, args map[string]string) (*NodeBalancerConfigResponse, error) {
	u := &url.Values{}
	u.Add("ConfigID", strconv.Itoa(configId))

	// add optional parameters
	processOptionalArgs(args, u)
	v := NodeBalancerConfigResponse{}
	if err := t.client.do("nodebalancer.config.update", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerConfigId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Delete Config
func (t *NodeBalancerConfigService) Delete(nodeBalancerId int, configId int) (*NodeBalancerConfigResponse, error) {
	u := &url.Values{}
	u.Add("NodeBalancerID", strconv.Itoa(nodeBalancerId))
	u.Add("ConfigID", strconv.Itoa(configId))
	v := NodeBalancerConfigResponse{}
	if err := t.client.do("nodebalancer.config.delete", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerConfigId); err != nil {
		return nil, err
	}
	return &v, nil
}

