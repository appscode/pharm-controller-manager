package scaleway

import (
	"context"

	"github.com/pharmer/cloud-controller-manager/cloud"
	scw "github.com/scaleway/scaleway-cli/pkg/api"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

type zones struct {
	client *scw.ScalewayAPI
	region string
}

func newZones(client *scw.ScalewayAPI, region string) cloudprovider.Zones {
	return &zones{client, region}
}

func (z zones) GetZone(_ context.Context) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, cloud.ErrNotImplemented
}

func (z zones) GetZoneByProviderID(_ context.Context, providerID string) (cloudprovider.Zone, error) {
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	server, err := serverByID(z.client, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: server.Location.ZoneID}, nil
}

func (z zones) GetZoneByNodeName(_ context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	server, err := serverByName(z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: server.Location.ZoneID}, nil
}
