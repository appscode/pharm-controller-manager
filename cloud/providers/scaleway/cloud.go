package scaleway

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
	scw "github.com/scaleway/scaleway-cli/pkg/api"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	ProviderName = "scaleway"
)

type Credential struct {
	Organization string `json:"organization" yaml:"organization"`
	Token        string `json:"token" yaml:"token"`
	Region       string `json:"region" yaml:"region"`
}
type Cloud struct {
	client        *scw.ScalewayAPI
	instances     cloudprovider.Instances
	zones         cloudprovider.Zones
	loadbalancers cloudprovider.LoadBalancer
}

func init() {
	cloudprovider.RegisterCloudProvider(
		ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			return newCloud(config)
		})
}

func newCloud(config io.Reader) (*Cloud, error) {
	cred := &Credential{}
	contents, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(contents))

	err = yaml.Unmarshal(contents, cred)
	if err != nil {
		return nil, err
	}
	client, err := scw.NewScalewayAPI(cred.Organization, cred.Token, "pharmer", cred.Region)
	if err != nil {
		return nil, err
	}

	return &Cloud{
		client:        client,
		instances:     newInstances(client),
		zones:         newZones(client, cred.Region),
		loadbalancers: newLoadbalancers(client),
	}, nil
}

func (c *Cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {
}

func (c *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return c.loadbalancers, true
}

func (c *Cloud) Instances() (cloudprovider.Instances, bool) {
	return c.instances, true
}

func (c *Cloud) Zones() (cloudprovider.Zones, bool) {
	return c.zones, true
}

func (c *Cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *Cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (c *Cloud) ProviderName() string {
	return ProviderName
}

func (c *Cloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

func (c *Cloud) HasClusterID() bool {
	return false
}
