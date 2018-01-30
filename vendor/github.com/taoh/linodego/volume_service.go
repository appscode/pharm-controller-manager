package linodego

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// LinodeVolumeService is an interface generated for "github.com/appscode/linodego".LinodeVolumeService.
type LinodeVolumeInterface interface {
	Clone(int, string) (*LinodeVolumeResponse, error)
	Create(int, string, map[string]string) (*LinodeVolumeResponse, error)
	Delete(int) (*LinodeVolumeResponse, error)
	List(int) (*LinodeVolumeListResponse, error)
	Update(int, map[string]string) (*LinodeVolumeResponse, error)
}

// Linode Volume Service
type LinodeVolumeService struct {
	client *Client
}

var _ LinodeVolumeInterface = &LinodeVolumeService{}

// Response for volume.list API
type LinodeVolumeListResponse struct {
	Response
	Volume []Volume
}

// Response for general config APIs
type LinodeVolumeResponse struct {
	Response
	VolumeId StackScriptId
}

// List all volumes. If volumeId is greater than 0, limit the results to given volume.
func (t *LinodeVolumeService) List(volumeId int) (*LinodeVolumeListResponse, error) {
	u := &url.Values{}
	if volumeId > 0 {
		u.Add("VolumeID", strconv.Itoa(volumeId))
	}
	v := LinodeVolumeListResponse{}
	if err := t.client.do("volume.list", u, &v.Response); err != nil {
		return nil, err
	}

	v.Volume = make([]Volume, 5)
	if err := json.Unmarshal(v.RawData, &v.Volume); err != nil {
		return nil, err
	}
	return &v, nil
}

// Create Volume
func (t *LinodeVolumeService) Create(size int, label string, args map[string]string) (*LinodeVolumeResponse, error) {
	u := &url.Values{}
	u.Add("Size", strconv.Itoa(size))
	u.Add("Label", label)
	// add optional parameters
	processOptionalArgs(args, u)
	v := LinodeVolumeResponse{}
	if err := t.client.do("volume.create", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.VolumeId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update Volume. See https://www.linode.com/api/volume/volume.update for allowed arguments.
func (t *LinodeVolumeService) Update(volumeId int, args map[string]string) (*LinodeVolumeResponse, error) {
	u := &url.Values{}
	u.Add("VolumeID", strconv.Itoa(volumeId))
	// add optional parameters
	processOptionalArgs(args, u)
	v := LinodeVolumeResponse{}
	if err := t.client.do("volume.update", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.VolumeId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Delete Volume
func (t *LinodeVolumeService) Delete(volumeId int) (*LinodeVolumeResponse, error) {
	u := &url.Values{}
	u.Add("VolumeID", strconv.Itoa(volumeId))
	v := LinodeVolumeResponse{}
	if err := t.client.do("volume.delete", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.VolumeId); err != nil {
		return nil, err
	}
	return &v, nil
}

// Clone Volume
func (t *LinodeVolumeService) Clone(cloneFromId int, label string) (*LinodeVolumeResponse, error) {
	u := &url.Values{}
	u.Add("CloneFromID", strconv.Itoa(cloneFromId))
	u.Add("Label", label)
	v := LinodeVolumeResponse{}
	if err := t.client.do("volume.clone", u, &v.Response); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.RawData, &v.VolumeId); err != nil {
		return nil, err
	}
	return &v, nil
}
