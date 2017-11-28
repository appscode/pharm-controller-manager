package lightsail

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/appscode/go/types"
	"github.com/appscode/pharm-controller-manager/cloud"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type instances struct {
	client *lightsail.Lightsail
}

func newInstances(client *lightsail.Lightsail) cloudprovider.Instances {
	return &instances{client}
}

func (i *instances) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	instance, err := instanceByName(i.client, name)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(instance)
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	id, err := instanceIDFromProviderID(providerID)
	if err != nil {
		return nil, err
	}
	instance, err := instanceByID(i.client, id)
	if err != nil {
		return nil, err
	}

	return nodeAddresses(instance)
}

func nodeAddresses(instance *lightsail.Instance) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: String(instance.Name)})

	if *instance.PrivateIpAddress == "" {
		return nil, fmt.Errorf("could not get private ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: String(instance.PrivateIpAddress)})

	if *instance.PublicIpAddress == "" {
		return nil, fmt.Errorf("could not get public ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: String(instance.PublicIpAddress)})

	return addresses, nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	instance, err := instanceByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return *instance.Arn, nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	instance, err := instanceByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return *instance.BundleId, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	id, err := instanceIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}
	instance, err := instanceByID(i.client, id)
	if err != nil {
		return "", err
	}
	return *instance.BundleId, nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *instances) InstanceExistsByProviderID(providerID string) (bool, error) {
	id, err := instanceIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}

	_, err = instanceByID(i.client, id)
	if err == nil {
		return true, nil
	}

	return false, nil
}

func instanceByName(client *lightsail.Lightsail, nodeName types.NodeName) (*lightsail.Instance, error) {
	hosts, err := allInstanceList(client)
	if err != nil {
		return nil, err
	}

	for _, host := range hosts {
		if getNodeName(*host.PrivateIpAddress) == nodeName {
			return host, nil
		}
	}

	return nil, cloudprovider.InstanceNotFound

}

func instanceByID(client *lightsail.Lightsail, providerId string) (*lightsail.Instance, error) {
	hosts, err := allInstanceList(client)
	if err != nil {
		return nil, err
	}

	for _, host := range hosts {
		if *host.Arn == providerId {
			return host, nil
		}
	}

	return nil, cloudprovider.InstanceNotFound

}

func getNodeName(ip string) types.NodeName {
	address := strings.Replace(ip, ".", "-", -1)
	name := "ip-" + address
	return types.NodeName(name)
}

// instanceIDFromProviderID returns a server's ID from providerID.
//
// The providerID spec should be retrievable from the Kubernetes
// node object. The expected format is: lightsail://server-id

func instanceIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be empty string")
	}

	split := strings.Split(providerID, "/")
	if len(split) != 3 {
		return "", fmt.Errorf("unexpected providerID format: %s, format should be: lightsail://12345", providerID)
	}

	// since split[0] is actually "vultr:"
	if strings.TrimSuffix(split[0], ":") != ProviderName {
		return "", fmt.Errorf("provider name from providerID should be vultr: %s", providerID)
	}

	return split[2], nil
}
