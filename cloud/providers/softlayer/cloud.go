package softlayer

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	ProviderName = "softlayer"
)

type Credential struct {
	UserName string `json:"username" yaml:"username"`
	ApiKey   string `json:"apiKey" yaml:"apiKey"`
	Zone     string `json:"zone" yaml:"zone"`
}

type Cloud struct {
	virtualServiceClient services.Virtual_Guest
	accountServiceClient services.Account

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

	sess := session.New(cred.UserName, cred.ApiKey)
	virtualServiceClient := services.GetVirtualGuestService(sess)
	accountServiceClient := services.GetAccountService(sess)

	return &Cloud{
		virtualServiceClient: virtualServiceClient,
		accountServiceClient: accountServiceClient,

		instances:     newInstances(virtualServiceClient, accountServiceClient),
		zones:         newZones(virtualServiceClient, accountServiceClient, cred.Zone),
		loadbalancers: newLoadbalancers(virtualServiceClient, accountServiceClient),
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
	return true
}
