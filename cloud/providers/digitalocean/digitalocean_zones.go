package digitalocean

import (
	"io/ioutil"
	"net/http"
	"sync"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

var faultMutex = &sync.Mutex{}

const instanceInfoURL = "http://169.254.169.254/metadata/v1"

func (d *DO) GetZone() (cloudprovider.Zone, error) {
	faultMutex.Lock()
	region, err := fetchRegion()
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	zone := cloudprovider.Zone{
		FailureDomain: region,
		Region:        region,
	}
	faultMutex.Unlock()
	return zone, nil

}

func (d *DO) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, nil
}

func (d *DO) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, nil
}

func fetchRegion() (string, error) {
	resp, err := http.Get(instanceInfoURL + "/region")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
