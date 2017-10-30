package softlayer

import (
	"github.com/softlayer/softlayer-go/services"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type zones struct {
	virtualServiceClient services.Virtual_Guest
	accountServiceClient services.Account

	region string
}

func newZones(virtualServiceClient services.Virtual_Guest,
	accountServiceClient services.Account, region string) cloudprovider.Zones {
	return &zones{virtualServiceClient: virtualServiceClient,
		accountServiceClient: accountServiceClient, region: region}
}

func (z zones) GetZone() (cloudprovider.Zone, error) {
	return cloudprovider.Zone{Region: z.region}, nil
}

func (z zones) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	id, err := guestIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	vGuest, err := guestByID(z.virtualServiceClient, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: *vGuest.Datacenter.Name}, nil

}

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	vGuest, err := guestByName(z.accountServiceClient, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: *vGuest.Datacenter.Name}, nil
}
