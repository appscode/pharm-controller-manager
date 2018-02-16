package packet

import (
	"fmt"
	"strings"

	"github.com/packethost/packngo"
	"github.com/pharmer/cloud-controller-manager/cloud"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type instances struct {
	client  *packngo.Client
	project string
}

func newInstances(client *packngo.Client, projectID string) cloudprovider.Instances {
	return &instances{client, projectID}
}

func (i *instances) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	device, err := deviceByName(i.client, i.project, name)
	if err != nil {
		return nil, err
	}
	return i.nodeAddresses(device)
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	id, err := deviceIDFromProviderID(providerID)
	if err != nil {
		return nil, err
	}
	device, err := deviceByID(i.client, id)
	if err != nil {
		return nil, err
	}

	return i.nodeAddresses(device)
}

func (i *instances) nodeAddresses(device *packngo.Device) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: device.Hostname})

	host, _, err := i.client.Devices.Get(device.ID)
	if err != nil {
		return nil, err
	}
	var privateIP, publicIP string

	for _, addr := range host.Network {
		if addr.AddressFamily == 4 {
			if addr.Public {
				publicIP = addr.Address
			} else {
				privateIP = addr.Address
			}
		}
	}
	if privateIP == "" {
		return nil, fmt.Errorf("could not get private ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: privateIP})

	if publicIP == "" {
		return nil, fmt.Errorf("could not get public ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: publicIP})

	return addresses, nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	device, err := deviceByName(i.client, i.project, nodeName)
	if err != nil {
		return "", err
	}
	return device.ID, nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	device, err := deviceByName(i.client, i.project, nodeName)
	if err != nil {
		return "", err
	}
	return device.Plan.Slug, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	id, err := deviceIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}
	device, err := deviceByID(i.client, id)
	if err != nil {
		return "", err
	}
	return device.Plan.Slug, nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *instances) InstanceExistsByProviderID(providerID string) (bool, error) {
	id, err := deviceIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}
	_, err = deviceByID(i.client, id)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func deviceByID(client *packngo.Client, id string) (*packngo.Device, error) {
	device, _, err := client.Devices.Get(id)
	return device, err
}

func deviceByName(client *packngo.Client, projectID string, nodeName types.NodeName) (*packngo.Device, error) {
	devices, _, err := client.Devices.List(projectID)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.Hostname == string(nodeName) {
			return &device, nil
		}
	}
	return nil, cloudprovider.InstanceNotFound
}

// deviceIDFromProviderID returns a device's ID from providerID.
//
// The providerID spec should be retrievable from the Kubernetes
// node object. The expected format is: packet://device-id

func deviceIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be empty string")
	}

	split := strings.Split(providerID, "/")
	if len(split) != 3 {
		return "", fmt.Errorf("unexpected providerID format: %s, format should be: packet://12345", providerID)
	}

	// since split[0] is actually "packet:"
	if strings.TrimSuffix(split[0], ":") != ProviderName {
		return "", fmt.Errorf("provider name from providerID should be packet: %s", providerID)
	}

	return split[2], nil
}
