package digitalocean

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	ProviderName = "digitalocean"
)

type Config struct {
	Token string `json:"token" yaml:"token"`
}
type DO struct {
	Config
	Client *godo.Client
}

func init() {
	cloudprovider.RegisterCloudProvider(
		ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			return newDO(config)
		})
}

func newDO(config io.Reader) (*DO, error) {
	var do DO
	contents, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(contents))

	err = json.Unmarshal(contents, &do)
	if err != nil {
		return nil, err
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: do.Token,
	}))
	do.Client = godo.NewClient(oauthClient)
	return &do, nil
}

func (d *DO) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (d *DO) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (d *DO) Instances() (cloudprovider.Instances, bool) {
	return d, true
}

func (d *DO) Zones() (cloudprovider.Zones, bool) {
	return d, true
}

func (d *DO) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (d *DO) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (d *DO) ProviderName() string {
	return ProviderName
}

func (d *DO) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nameservers, searches
}

func (d *DO) HasClusterID() bool {
	return false
}
