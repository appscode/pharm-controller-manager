package vultr

import (
	"io/ioutil"
	"net/http"
	"strconv"

	gv "github.com/JamesClonk/vultr/lib"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const (
	serverIDURL   = "http://169.254.169.254/latest/meta-data/SUBID"
	serverListURL = "https://api.vultr.com/v1/server/list"
)

type zones struct {
	client *gv.Client
}

func newZones(client *gv.Client) cloudprovider.Zones {
	return zones{client}
}

func (z zones) GetZone() (cloudprovider.Zone, error) {
	subid, err := fetchServerID()
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	server, err := serverByID(z.client, subid)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: strconv.Itoa(server.RegionID)}, nil
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

	return cloudprovider.Zone{Region: strconv.Itoa(server.RegionID)}, nil
}

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	server, err := serverByName(z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{Region: strconv.Itoa(server.RegionID)}, nil
}

func fetchRegion(token string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", serverListURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("API-Key", token)

	subid, err := fetchServerID()
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("SUBID", subid)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func fetchServerID() (string, error) {
	resp, err := http.Get(serverIDURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
