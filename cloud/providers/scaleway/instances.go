package scaleway

import (
	"github.com/appscode/pharm-controller-manager/cloud"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *Cloud) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	return nil, cloud.ErrNotImplemented
}

func (c *Cloud) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	return nil, cloud.ErrNotImplemented
}

func (c *Cloud) ExternalID(nodeName types.NodeName) (string, error) {
	return "", cloud.ErrNotImplemented
}

func (c *Cloud) InstanceID(nodeName types.NodeName) (string, error) {
	return "", cloud.ErrNotImplemented
}

func (c *Cloud) InstanceType(nodeName types.NodeName) (string, error) {
	return "", cloud.ErrNotImplemented
}

func (c *Cloud) InstanceTypeByProviderID(providerID string) (string, error) {
	return "", cloud.ErrNotImplemented
}

func (c *Cloud) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (c *Cloud) CurrentNodeName(hostname string) (types.NodeName, error) {
	return "", cloud.ErrNotImplemented
}

func (c *Cloud) InstanceExistsByProviderID(providerID string) (bool, error) {
	return false, cloud.ErrNotImplemented
}
