package softlayer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"pharmer.dev/cloud-controller-manager/cloud"
)

type instances struct {
	virtualServiceClient services.Virtual_Guest
	accountServiceClient services.Account
}

func newInstances(virtualServiceClient services.Virtual_Guest,
	accountServiceClient services.Account) cloudprovider.Instances {
	return &instances{virtualServiceClient: virtualServiceClient,
		accountServiceClient: accountServiceClient}
}

func (i *instances) NodeAddresses(_ context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	vGuest, err := guestByName(i.accountServiceClient, name)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(i.virtualServiceClient, vGuest)
}

func (i *instances) NodeAddressesByProviderID(_ context.Context, providerID string) ([]v1.NodeAddress, error) {
	id, err := guestIDFromProviderID(providerID)
	if err != nil {
		return nil, err
	}

	vGuest, err := guestByID(i.virtualServiceClient, id)
	if err != nil {
		return nil, err
	}

	return nodeAddresses(i.virtualServiceClient, vGuest)
}

func nodeAddresses(virtualServiceClient services.Virtual_Guest, vGuest datatypes.Virtual_Guest) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: *vGuest.Hostname})

	bluemix := virtualServiceClient.Id(*vGuest.Id)
	privateIP, err := bluemix.GetPrimaryBackendIpAddress()
	if err != nil {
		return nil, fmt.Errorf("could not get private ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: privateIP})

	publicIP, err := bluemix.GetPrimaryIpAddress()
	if err != nil {
		return nil, fmt.Errorf("could not get public ip")
	}

	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: publicIP})

	return addresses, nil
}

func (i *instances) ExternalID(ctx context.Context, nodeName types.NodeName) (string, error) {
	return i.InstanceID(ctx, nodeName)
}

func (i *instances) InstanceID(_ context.Context, nodeName types.NodeName) (string, error) {
	vGuest, err := guestByName(i.accountServiceClient, nodeName)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(*vGuest.Id), nil
}

func (i *instances) InstanceType(_ context.Context, nodeName types.NodeName) (string, error) {
	vGuest, err := guestByName(i.accountServiceClient, nodeName)
	if err != nil {
		return "", err
	}

	return guestInstanceType(vGuest)
}

func (i *instances) InstanceTypeByProviderID(_ context.Context, providerID string) (string, error) {
	id, err := guestIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}

	vGuest, err := guestByID(i.virtualServiceClient, id)
	if err != nil {
		return "", err
	}

	return guestInstanceType(vGuest)
}

func (i *instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, cloudprovider.NotImplemented
}

func guestInstanceType(vGuest datatypes.Virtual_Guest) (string, error) {
	cpu := *vGuest.StartCpus
	ram := (*vGuest.MaxMemory) / 1024
	sku := fmt.Sprintf("%vc%vm", cpu, ram)
	return sku, nil
}

func (i *instances) AddSSHKeyToAllInstances(_ context.Context, user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (i *instances) CurrentNodeName(_ context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *instances) InstanceExistsByProviderID(_ context.Context, providerID string) (bool, error) {
	id, err := guestIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}

	_, err = guestByID(i.virtualServiceClient, id)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func guestByID(virtualServiceClient services.Virtual_Guest, id string) (datatypes.Virtual_Guest, error) {
	guestID, err := strconv.Atoi(id)
	if err != nil {
		return datatypes.Virtual_Guest{}, err
	}

	vGuest, err := virtualServiceClient.Id(guestID).GetObject()
	if err != nil {
		return datatypes.Virtual_Guest{}, err
	}
	return vGuest, nil
}

func guestByName(accountServiceClient services.Account, nodeName types.NodeName) (datatypes.Virtual_Guest, error) {
	guests, err := accountServiceClient.GetVirtualGuests()
	if err != nil {
		return datatypes.Virtual_Guest{}, err
	}
	for _, guest := range guests {
		if *guest.Hostname == string(nodeName) {
			return guest, err
		}
	}
	return datatypes.Virtual_Guest{}, cloudprovider.InstanceNotFound
}

// serverIDFromProviderID returns a server's ID from providerID.
//
// The providerID spec should be retrievable from the Kubernetes
// node object. The expected format is: softlayer://server-id

func guestIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be empty string")
	}

	split := strings.Split(providerID, "/")
	if len(split) != 3 {
		return "", fmt.Errorf("unexpected providerID format: %s, format should be: softlayer://12345", providerID)
	}

	// since split[0] is actually "softlayer:"
	if strings.TrimSuffix(split[0], ":") != ProviderName {
		return "", fmt.Errorf("provider name from providerID should be softlayer: %s", providerID)
	}

	return split[2], nil
}
