package lightsail

import (
	"fmt"

	. "github.com/appscode/go/types"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const (
	metadataURL = "http://169.254.169.254/latest/meta-data/"
)

type zones struct {
	client *lightsail.Lightsail
}

func newZones(client *lightsail.Lightsail) cloudprovider.Zones {
	return zones{client}
}

func (z zones) GetZone() (cloudprovider.Zone, error) {
	return getZone()
}

func (z zones) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	fmt.Println(providerID, "-----------------------")
	return getZone()
}

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	instance, err := instanceByName(z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: String(instance.Location.RegionName), FailureDomain: String(instance.Location.AvailabilityZone)}, nil
}

func getZone() (cloudprovider.Zone, error) {
	zone, err := getAvailabilityZone()
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	region, err := azToRegion(zone)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: region, FailureDomain: zone}, nil
}

// Derives the region from a valid az name.
// Returns an error if the az is known invalid (empty)
func azToRegion(az string) (string, error) {
	if len(az) < 1 {
		return "", fmt.Errorf("invalid (empty) AZ")
	}
	region := az[:len(az)-1]
	return region, nil
}

func getAvailabilityZone() (string, error) {
	zone := "placement/availability-zone"
	return GetMetadata(zone)
}
