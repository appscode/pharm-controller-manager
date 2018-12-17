package packet

import (
	"context"

	"github.com/packethost/packngo"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

type zones struct {
	client  *packngo.Client
	project string
	zone    string
}

func newZones(client *packngo.Client, projectID, zone string) cloudprovider.Zones {
	return zones{client, projectID, zone}
}

func (z zones) GetZone(_ context.Context) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{Region: z.zone}, nil
}

func (z zones) GetZoneByProviderID(_ context.Context, providerID string) (cloudprovider.Zone, error) {
	id, err := deviceIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	device, err := deviceByID(z.client, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: device.Facility.ID}, nil

}

func (z zones) GetZoneByNodeName(_ context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	device, err := deviceByName(z.client, z.project, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: device.Facility.ID}, nil
}
