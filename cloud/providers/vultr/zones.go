package vultr

import (
	"github.com/appscode/pharm-controller-manager/cloud"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
	gv "github.com/JamesClonk/vultr/lib"
	"strconv"
)

type zones struct {
	client *gv.Client
}

func newZones(client *gv.Client) cloudprovider.Zones {
	return zones{client}
}


func (z zones) GetZone() (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, cloud.ErrNotImplemented
}

func (z zones) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	server, err := serverByID(z.client, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: strconv.Itoa(server.RegionID)}, cloud.ErrNotImplemented
}

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	server, err := serverByName(z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: strconv.Itoa(server.RegionID)}, cloud.ErrNotImplemented
}
