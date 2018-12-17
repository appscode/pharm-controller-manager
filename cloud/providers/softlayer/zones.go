package softlayer

import (
	"context"
	"strconv"

	"github.com/softlayer/softlayer-go/services"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

type zones struct {
	virtualServiceClient services.Virtual_Guest
	accountServiceClient services.Account

	zone string
}

func newZones(virtualServiceClient services.Virtual_Guest,
	accountServiceClient services.Account, region string) cloudprovider.Zones {
	return &zones{virtualServiceClient: virtualServiceClient,
		accountServiceClient: accountServiceClient, zone: region}
}

func (z zones) GetZone(_ context.Context) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{Region: z.zone}, nil
}

func (z zones) GetZoneByProviderID(_ context.Context, providerID string) (cloudprovider.Zone, error) {
	id, err := guestIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	location, err := fetchDatacenterLocation(z.virtualServiceClient, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: location}, nil

}

func (z zones) GetZoneByNodeName(_ context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	vGuest, err := guestByName(z.accountServiceClient, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	location, err := fetchDatacenterLocation(z.virtualServiceClient, strconv.Itoa(*vGuest.Id))
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: location}, nil
}

func fetchDatacenterLocation(virtualServiceClient services.Virtual_Guest, id string) (string, error) {
	guestID, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	datacenter, err := virtualServiceClient.Id(guestID).GetDatacenter()
	if err != nil {
		return "", err
	}
	return *datacenter.Name, nil
}
