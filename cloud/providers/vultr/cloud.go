package vultr

import (
	"io"
	"io/ioutil"

	gv "github.com/JamesClonk/vultr/lib"
	"github.com/ghodss/yaml"
	cloudprovider "k8s.io/cloud-provider"
)

const (
	ProviderName = "vultr"
)

type tokenSource struct {
	Token string `json:"token" yaml:"token"`
}

type Cloud struct {
	client        *gv.Client
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

func newCloud(config io.Reader) (cloudprovider.Interface, error) {
	tokenSource := &tokenSource{}
	contents, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(contents, tokenSource)
	if err != nil {
		return nil, err
	}

	vultrClient := gv.NewClient(tokenSource.Token, &gv.Options{})
	return &Cloud{
		client:        vultrClient,
		instances:     newInstances(vultrClient),
		zones:         newZones(vultrClient),
		loadbalancers: newLoadbalancers(vultrClient),
	}, nil
}

func (c *Cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
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
	return true
}
