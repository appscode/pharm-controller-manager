package scaleway

import (
	"github.com/appscode/pharm-controller-manager/cloud"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func (c *Cloud) GetZone() (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, cloud.ErrNotImplemented
}

func (c *Cloud) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, cloud.ErrNotImplemented
}

func (c *Cloud) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, cloud.ErrNotImplemented
}
