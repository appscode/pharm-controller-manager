package linode

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/appscode/pharm-controller-manager/cloud"
	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	ProviderName = "linode"
)

type Config struct {
	Token string `json:"token" yaml:"token"`
}
type Cloud struct {
	Config
}

func init() {
	cloudprovider.RegisterCloudProvider(
		ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			return newCloud(config)
		})
}

func newCloud(config io.Reader) (*Cloud, error) {
	var do Cloud
	contents, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(contents))

	err = yaml.Unmarshal(contents, &do)
	if err != nil {
		return nil, err
	}

	return nil, cloud.ErrNotImplemented
}

func (c *Cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (c *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (c *Cloud) Instances() (cloudprovider.Instances, bool) {
	return c, true
}

func (c *Cloud) Zones() (cloudprovider.Zones, bool) {
	return c, true
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
	return nameservers, searches
}

func (c *Cloud) HasClusterID() bool {
	return false
}
