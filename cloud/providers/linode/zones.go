package linode

import (
	"strconv"

	"github.com/taoh/linodego"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type zones struct {
	client *linodego.Client
	zone   string
}

func newZones(client *linodego.Client, zone string) cloudprovider.Zones {
	return zones{client, zone}
}

func (z zones) GetZone() (cloudprovider.Zone, error) {
	return cloudprovider.Zone{Region: z.zone}, nil
}

func (z zones) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	linode, err := linodeByID(z.client, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: strconv.Itoa(linode.DataCenterId)}, nil
}

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	linode, err := linodeByName(z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: strconv.Itoa(linode.DataCenterId)}, nil
}
