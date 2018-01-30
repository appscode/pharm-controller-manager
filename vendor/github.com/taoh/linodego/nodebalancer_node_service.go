package linodego

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// NodeBalancerNodeService is an interface generated for "github.com/appscode/linodego".NodeBalancerNodeService.
type NodeBalancerNodeInterface interface {
	Create(int, string, string, map[string]string) (*NodeBalancerNodeResponse, error)
	Delete(int) (*NodeBalancerNodeResponse, error)
	List(int, int) (*NodeBalancerNodeListResponse, error)
	Update(int, map[string]string) (*NodeBalancerNodeResponse, error)
}

// NodeBalancer Config Service
type NodeBalancerNodeService struct {
	client *Client
}

var _ NodeBalancerNodeInterface = &NodeBalancerNodeService{}

// Response for nodebalancer.config.list API
type NodeBalancerNodeListResponse struct {
	Response
	NodeBalancerNodes []NodeBalancerNode
}

// Response for general config APIs
type NodeBalancerNodeResponse struct {
	Response
	NodeBalancerNodeId NodeBalancerNodeId
}

// Get Node List. If configId is greater than 0, limit results to given node.
func (t *NodeBalancerNodeService) List(configId int, nodeId int) (*NodeBalancerNodeListResponse, error) {
	u := &url.Values{}
	u.Add("ConfigID", strconv.Itoa(configId))
	if nodeId > 0 {
		u.Add("NodeID", strconv.Itoa(nodeId))
	}
	v := NodeBalancerNodeListResponse{}
	if err := t.client.do("nodebalancer.node.list", u, &v.Response); err != nil {
		return nil, err
	}

	v.NodeBalancerNodes = make([]NodeBalancerNode, 5)
	if err := json.Unmarshal(v.RawData, &v.NodeBalancerNodes); err != nil {
		return nil, err
	}
	return &v, nil
}

// Create Node
func (t *NodeBalancerNodeService) Create(configId int, label string, address string, args map[string]string) (*NodeBalancerNodeResponse, error) {
	u := &url.Values{}
	u.Add("ConfigID", strconv.Itoa(configId))
	u.Add("Label", label)
	u.Add("Address", address)
	// add optional parameters
	processOptionalArgs(args, u)
	v := NodeBalancerNodeResponse{}
	if err := t.client.do("nodebalancer.node.create", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerNodeId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update Node. See https://www.linode.com/api/linode/linode.config.update for allowed arguments.
func (t *NodeBalancerNodeService) Update(NodeId int, args map[string]string) (*NodeBalancerNodeResponse, error) {
	u := &url.Values{}
	u.Add("NodeID", strconv.Itoa(NodeId))

	// add optional parameters
	processOptionalArgs(args, u)
	v := NodeBalancerNodeResponse{}
	if err := t.client.do("nodebalancer.node.update", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerNodeId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Delete Node
func (t *NodeBalancerNodeService) Delete(nodeId int) (*NodeBalancerNodeResponse, error) {
	u := &url.Values{}
	u.Add("NodeID", strconv.Itoa(nodeId))
	v := NodeBalancerNodeResponse{}
	if err := t.client.do("nodebalancer.node.delete", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.NodeBalancerNodeId); err != nil {
		return nil, err
	}
	return &v, nil
}
