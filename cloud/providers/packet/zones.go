package packet

import (
	"github.com/packethost/packngo"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type zones struct {
	client  *packngo.Client
	project string
	zone    string
}

func newZones(client *packngo.Client, projectID, zone string) cloudprovider.Zones {
	return zones{client, projectID, zone}
}

func (z zones) GetZone() (cloudprovider.Zone, error) {
	return cloudprovider.Zone{Region: z.zone}, nil
}

func (z zones) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
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

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	device, err := deviceByName(z.client, z.project, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: device.Facility.ID}, nil
}
