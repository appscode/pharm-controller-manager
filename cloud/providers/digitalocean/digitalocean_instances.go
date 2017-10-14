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

func (d *DO) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	droplet, err := d.getDroplet(name)
	if err != nil {
		return []v1.NodeAddress{}, err
	}

	return d.getNodeAddress(&droplet)
}

func (d *DO) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	return []v1.NodeAddress{}, errors.New("unimplemented:" + providerID)

	droplet, err := d.getDropletByID(providerID)
	if err != nil {
		return []v1.NodeAddress{}, err
	}

	return d.getNodeAddress(droplet)
}

func (d *DO) ExternalID(nodeName types.NodeName) (string, error) {
	droplet, err := d.getDroplet(nodeName)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(droplet.ID), nil
}

func (d *DO) InstanceID(nodeName types.NodeName) (string, error) {
	return d.ExternalID(nodeName)
}

func (d *DO) InstanceType(nodeName types.NodeName) (string, error) {
	droplet, err := d.getDroplet(nodeName)
	if err != nil {
		return "", err
	}
	return droplet.SizeSlug, nil
}

func (d *DO) InstanceTypeByProviderID(providerID string) (string, error) {
	return "", errors.New("unimplemented:" + providerID)
	droplet, err := d.getDropletByID(providerID)
	if err != nil {
		return "", err
	}
	return droplet.SizeSlug, nil
}

func (d *DO) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return errors.New("unimplemented")
}

func (d *DO) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (d *DO) InstanceExistsByProviderID(providerID string) (bool, error) {
	return false, nil
}

func (d *DO) getNodeAddress(droplet *godo.Droplet) ([]v1.NodeAddress, error) {
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

func (d *DO) getDroplet(name types.NodeName) (godo.Droplet, error) {
	droplets, _, err := d.Client.Droplets.List(context.TODO(), &godo.ListOptions{})
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

func (d *DO) getDropletByID(providerID string) (*godo.Droplet, error) {
	id, err := strconv.Atoi(providerID)
	if err != nil {
		return &godo.Droplet{}, err
	}
	droplet, _, err := d.Client.Droplets.Get(context.TODO(), id)
	if err != nil {
		return &godo.Droplet{}, err
	}
	return droplet, nil
}
