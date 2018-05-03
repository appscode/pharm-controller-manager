package lightsail

import (
	"context"
	"fmt"
	"strings"

	. "github.com/appscode/go/types"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/pharmer/cloud-controller-manager/cloud"
	"github.com/pkg/errors"
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

func (i *instances) NodeAddresses(_ context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	instance, err := instanceByName(i.client, name)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(instance)
}

func (i *instances) NodeAddressesByProviderID(_ context.Context, providerID string) ([]v1.NodeAddress, error) {
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

func (i *instances) ExternalID(_ context.Context, nodeName types.NodeName) (string, error) {
	return string(nodeName), nil
}

func (i *instances) InstanceID(_ context.Context, nodeName types.NodeName) (string, error) {
	return string(nodeName), nil
}

func (i *instances) InstanceType(_ context.Context, nodeName types.NodeName) (string, error) {
	instance, err := instanceByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return *instance.BundleId, nil
}

func (i *instances) InstanceTypeByProviderID(_ context.Context, providerID string) (string, error) {
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

func (i *instances) AddSSHKeyToAllInstances(_ context.Context, user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (i *instances) CurrentNodeName(_ context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *instances) InstanceExistsByProviderID(_ context.Context, providerID string) (bool, error) {
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
	host, err := client.GetInstance(&lightsail.GetInstanceInput{
		InstanceName: StringP(string(nodeName)),
	})
	if err != nil {
		return nil, err
	}
	if host.Instance != nil {
		return host.Instance, nil
	}
	return nil, cloudprovider.InstanceNotFound

}

func instanceByID(client *lightsail.Lightsail, id string) (*lightsail.Instance, error) {
	return instanceByName(client, types.NodeName(id))
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

	// since split[0] is actually "lightsail:"
	if strings.TrimSuffix(split[0], ":") != ProviderName {
		return "", fmt.Errorf("provider name from providerID should be vultr: %s", providerID)
	}

	return split[2], nil
}
