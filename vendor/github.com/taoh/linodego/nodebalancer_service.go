package linodego

import (
	"encoding/json"
	"net/url"
	"strconv"
)


// Nodebalancer Service
type NodeBalancerService struct {
	client *Client
}

// Response for nodebalancer.list API
type NodeBalancerListResponse struct {
	Response
	NodeBalancer []LinodeNodeBalancer
}


// Response for general config APIs
type NodeBalancerResponse struct {
	Response
	NodeBalancerId LinodeNodeBalancerId
}


// List all nodebalancers. If nodeBalancerId is greater than 0, limit the results to given nodebalancer.
func (t *NodeBalancerService) List(nodeBalancerId int) (*NodeBalancerListResponse, error) {
	u := &url.Values{}
	if nodeBalancerId > 0 {
		u.Add("NodeBalancerID", strconv.Itoa(nodeBalancerId))
	}
	v := NodeBalancerListResponse{}
	if err := t.client.do("nodebalancer.list", u, &v.Response); err != nil {
		return nil, err
	}

	v.NodeBalancer = make([]LinodeNodeBalancer, 5)
	if err := json.Unmarshal(v.RawData, &v.NodeBalancer); err != nil {
		return nil, err
	}
	return &v, nil
}

// Create NodeBalancer
func (t *NodeBalancerService) Create(datacenterId int,label string, args map[string]string) (*NodeBalancerResponse, error) {
	u := &url.Values{}
	u.Add("DatacenterID", strconv.Itoa(datacenterId))
	u.Add("Label", label)
	// add optional parameters
	processOptionalArgs(args, u)
	v := NodeBalancerResponse{}
	if err := t.client.do("nodebalancer.create", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update Nodebalancer. See https://www.linode.com/api/nodebalancer/nodebalancer.update for allowed arguments.
func (t *NodeBalancerService) Update(nodeBalancerId int, args map[string]string) (*NodeBalancerResponse, error) {
	u := &url.Values{}
	u.Add("NodeBalancerID", strconv.Itoa(nodeBalancerId))
	// add optional parameters
	processOptionalArgs(args, u)
	v := NodeBalancerResponse{}
	if err := t.client.do("nodebalancer.update", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Delete Nodebalancer
func (t *NodeBalancerService) Delete(nodeBalancerId int) (*NodeBalancerResponse, error) {
	u := &url.Values{}
	u.Add("NodeBalancerID", strconv.Itoa(nodeBalancerId))
	v := NodeBalancerResponse{}
	if err := t.client.do("volume.delete", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerId); err != nil {
		return nil, err
	}
	return &v, nil
}

