package scaleway

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pharmer/cloud-controller-manager/cloud"
	scw "github.com/scaleway/scaleway-cli/pkg/api"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type instances struct {
	client *scw.ScalewayAPI
}

func newInstances(client *scw.ScalewayAPI) cloudprovider.Instances {
	return &instances{client}
}

func (i *instances) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	server, err := serverByName(i.client, name)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(server)
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return nil, err
	}
	server, err := serverByID(i.client, id)
	if err != nil {
		return nil, err
	}

	return nodeAddresses(server)
}

func nodeAddresses(server *scw.ScalewayServer) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: server.Name})

	if server.PrivateIP == "" {
		return nil, fmt.Errorf("could not get private ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: server.PrivateIP})

	if server.PublicAddress.IP == "" {
		return nil, fmt.Errorf("could not get public ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: server.PublicAddress.IP})

	return addresses, nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	server, err := serverByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return server.Identifier, nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	server, err := serverByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return server.CommercialType, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}
	server, err := serverByID(i.client, id)
	if err != nil {
		return "", err
	}
	return server.CommercialType, nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *instances) InstanceExistsByProviderID(providerID string) (bool, error) {
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}
	_, err = serverByID(i.client, id)
	if err == nil {
		return true, nil
	}

	return false, nil
}

func serverByID(client *scw.ScalewayAPI, id string) (*scw.ScalewayServer, error) {
	return client.GetServer(id)
}

func serverByName(client *scw.ScalewayAPI, nodeName types.NodeName) (*scw.ScalewayServer, error) {
	servers, err := client.GetServers(true, 0)
	if err != nil {
		return nil, err
	}

	for _, server := range *servers {
		if strings.ToLower(server.Name) == string(nodeName) {
			return &server, nil
		}
	}
	return nil, cloudprovider.InstanceNotFound
}

// serverIDFromProviderID returns a server's ID from providerID.
//
// The providerID spec should be retrievable from the Kubernetes
// node object. The expected format is: scaleway://server-id

func serverIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be empty string")
	}

	split := strings.Split(providerID, "/")
	if len(split) != 3 {
		return "", fmt.Errorf("unexpected providerID format: %s, format should be: scaleway://12345", providerID)
	}

	// since split[0] is actually "vultr:"
	if strings.TrimSuffix(split[0], ":") != ProviderName {
		return "", fmt.Errorf("provider name from providerID should be scaleway: %s", providerID)
	}

	return split[2], nil
}
