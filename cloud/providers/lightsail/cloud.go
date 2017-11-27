package lightsail

import (
	"io"
	"io/ioutil"
	//. "github.com/appscode/go/types"
	_aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	//"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"net/http"
)

const (
	ProviderName = "lightsail"
)

type tokenSource struct {
	AccessKeyID     string `json:"accessKeyID" yaml:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`
}

type Cloud struct {
	client        *lightsail.Lightsail
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
	zone, err := getZone()
	if err != nil {
		return nil, err
	}

	conf := &_aws.Config{
		Region:      &zone.Region,
		Credentials: credentials.NewStaticCredentials(tokenSource.AccessKeyID, tokenSource.SecretAccessKey, ""),
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}
	lightsailClient := lightsail.New(sess)

	return &Cloud{
		client:        lightsailClient,
		instances:     newInstances(lightsailClient),
		zones:         newZones(lightsailClient),
		loadbalancers: newLoadbalancers(lightsailClient),
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

func GetMetadata(path string) (string, error) {
	resp, err := http.Get(metadataURL + path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
