package digitalocean

import (
	"context"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func (c *Cloud) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	droplet, err := c.getDroplet(name)
	if err != nil {
		return []v1.NodeAddress{}, err
	}

	return c.getNodeAddress(&droplet)
}

func (c *Cloud) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	return []v1.NodeAddress{}, errors.New("unimplemented:" + providerID)

	droplet, err := c.getDropletByID(providerID)
	if err != nil {
		return []v1.NodeAddress{}, err
	}

	return c.getNodeAddress(droplet)
}

func (c *Cloud) ExternalID(nodeName types.NodeName) (string, error) {
	droplet, err := c.getDroplet(nodeName)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(droplet.ID), nil
}

func (c *Cloud) InstanceID(nodeName types.NodeName) (string, error) {
	return c.ExternalID(nodeName)
}

func (c *Cloud) InstanceType(nodeName types.NodeName) (string, error) {
	droplet, err := c.getDroplet(nodeName)
	if err != nil {
		return "", err
	}
	return droplet.SizeSlug, nil
}

func (c *Cloud) InstanceTypeByProviderID(providerID string) (string, error) {
	return "", errors.New("unimplemented:" + providerID)
	droplet, err := c.getDropletByID(providerID)
	if err != nil {
		return "", err
	}
	return droplet.SizeSlug, nil
}

func (c *Cloud) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return errors.New("unimplemented")
}

func (c *Cloud) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (c *Cloud) InstanceExistsByProviderID(providerID string) (bool, error) {
	return false, nil
}

func (c *Cloud) getNodeAddress(droplet *godo.Droplet) ([]v1.NodeAddress, error) {
	privateIP, err := droplet.PrivateIPv4()
	if err != nil {
		return []v1.NodeAddress{}, err
	}

	publicIP, err := droplet.PublicIPv4()
	if err != nil {
		return []v1.NodeAddress{}, err
	}
	return []v1.NodeAddress{
		{Type: v1.NodeInternalIP, Address: privateIP},
		{Type: v1.NodeExternalIP, Address: publicIP},
	}, nil
}

func (c *Cloud) getDroplet(name types.NodeName) (godo.Droplet, error) {
	droplets, _, err := c.client.Droplets.List(context.TODO(), &godo.ListOptions{})
	if err != nil {
		return godo.Droplet{}, err
	}
	nodeName := string(name)
	for _, item := range droplets {
		if item.Name == nodeName {
			return item, nil
		}
	}
	return godo.Droplet{}, cloudprovider.InstanceNotFound
}

func (c *Cloud) getDropletByID(providerID string) (*godo.Droplet, error) {
	id, err := strconv.Atoi(providerID)
	if err != nil {
		return &godo.Droplet{}, err
	}
	droplet, _, err := c.client.Droplets.Get(context.TODO(), id)
	if err != nil {
		return &godo.Droplet{}, err
	}
	return droplet, nil
}
