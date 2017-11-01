package linodego

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Linode Config Service
type StackScriptService struct {
	client *Client
}

// Response for linode.config.list API
type StackScriptListResponse struct {
	Response
	StackScripts []StackScript
}

// Response for general config APIs
type StackScriptResponse struct {
	Response
	StackScriptId StackScriptId
}

// Get Config List. If scriptId is greater than 0, limit results to given config.
func (t *StackScriptService) List(scriptId int) (*StackScriptListResponse, error) {
	u := &url.Values{}
	if scriptId > 0 {
		u.Add("StackScriptID", strconv.Itoa(scriptId))
	}
	v := StackScriptListResponse{}
	if err := t.client.do("stackscript.list", u, &v.Response); err != nil {
		return nil, err
	}

	v.StackScripts = make([]StackScript, 0)
	if err := json.Unmarshal(v.RawData, &v.StackScripts); err != nil {
		return nil, err
	}
	return &v, nil
}

// Create Config
func (t *StackScriptService) Create(label, distributionIDList, script string, args map[string]string) (*StackScriptResponse, error) {
	u := &url.Values{}
	u.Add("Label", label)
	u.Add("DistributionIDList", distributionIDList)
	u.Add("script", script)
	// add optional parameters
	processOptionalArgs(args, u)
	v := StackScriptResponse{}
	if err := t.client.do("stackscript.create", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.StackScriptId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update Config. See https://www.linode.com/api/stackscript/stackscript.update for allowed arguments.
func (t *StackScriptService) Update(scriptId int, args map[string]string) (*StackScriptResponse, error) {
	u := &url.Values{}
	u.Add("StackScriptID", strconv.Itoa(scriptId))
	// add optional parameters
	processOptionalArgs(args, u)
	v := StackScriptResponse{}
	if err := t.client.do("stackscript.update", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.StackScriptId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Delete Config
func (t *StackScriptService) Delete(scriptId int) (*StackScriptResponse, error) {
	u := &url.Values{}
	u.Add("StackScriptID", strconv.Itoa(scriptId))
	v := StackScriptResponse{}
	if err := t.client.do("stackscript.delete", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.StackScriptId); err != nil {
		return nil, err
	}
	return &v, nil
}
