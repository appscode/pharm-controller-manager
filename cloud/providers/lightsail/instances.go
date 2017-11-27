package lightsail

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/appscode/go/types"
	"github.com/appscode/pharm-controller-manager/cloud"
	_aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type awsInstance struct {
	// id in AWS
	awsID string

	// node name in k8s
	nodeName types.NodeName

	// availability zone the instance resides in
	availabilityZone string

	// ID of VPC the instance resides in
	vpcID string

	// ID of subnet the instance resides in
	subnetID string

	// instance type
	instanceType string
}

type instances struct {
	client *lightsail.Lightsail
}

func newInstances(client *lightsail.Lightsail) cloudprovider.Instances {
	return &instances{client}
}

func (i *instances) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	instance, err := instanceByName(i.client, name)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(instance)
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	fmt.Println(providerID, "**************8")
	return nil, nil
	/*id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return nil, err
	}
	//server, err := serverByID(i.client, id)
	if err != nil {
		return nil, err
	}

	return nodeAddresses(&server)*/
}

func nodeAddresses(instance *lightsail.Instance) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: String(instance.Name)})

	if *instance.PrivateIpAddress == "" {
		return nil, fmt.Errorf("could not get private ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: String(instance.PrivateIpAddress)})

	if *instance.PublicIpAddress == "" {
		return nil, fmt.Errorf("could not get public ip")
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: String(instance.PublicIpAddress)})

	return addresses, nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	instance, err := instanceByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return *instance.BlueprintId, nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	instance, err := instanceByName(i.client, nodeName)
	if err != nil {
		return "", err
	}
	return *instance.ResourceType, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	fmt.Println(providerID, "...............")
	/*	id, err := serverIDFromProviderID(providerID)
		if err != nil {
			return "", err
		}
		server, err := serverByID(i.client, id)
		if err != nil {
			return "", err
		}*/
	return "", nil
	//return strconv.Itoa(server.PlanID), nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return cloud.ErrNotImplemented
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i *instances) InstanceExistsByProviderID(providerID string) (bool, error) {
	//TODO(sanjid): check provider id here
	id, err := serverIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}
	fmt.Println(id, ",,,,,,,,,,,,,,,,")
	//_, err = insta(i.client, id)
	if err == nil {
		return true, nil
	}

	return false, nil
}

func instanceByName(client *lightsail.Lightsail, nodeName types.NodeName) (*lightsail.Instance, error) {
	host, err := client.GetInstance(&lightsail.GetInstanceInput{
		InstanceName: _aws.String(string(nodeName)),
	})

	if err != nil {
		return nil, err
	}
	return host.Instance, nil
}

// serverIDFromProviderID returns a server's ID from providerID.
//
// The providerID spec should be retrievable from the Kubernetes
// node object. The expected format is: lightsail://server-id

func serverIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be empty string")
	}

	split := strings.Split(providerID, "/")
	if len(split) != 3 {
		return "", fmt.Errorf("unexpected providerID format: %s, format should be: lightsail://12345", providerID)
	}

	// since split[0] is actually "vultr:"
	if strings.TrimSuffix(split[0], ":") != ProviderName {
		return "", fmt.Errorf("provider name from providerID should be vultr: %s", providerID)
	}

	return split[2], nil
}

func newAWSInstance(instance *ec2.Instance) *awsInstance {
	az := ""
	if instance.Placement != nil {
		az = _aws.StringValue(instance.Placement.AvailabilityZone)
	}
	self := &awsInstance{
		awsID:            _aws.StringValue(instance.InstanceId),
		nodeName:         mapInstanceToNodeName(instance),
		availabilityZone: az,
		instanceType:     _aws.StringValue(instance.InstanceType),
		vpcID:            _aws.StringValue(instance.VpcId),
		subnetID:         _aws.StringValue(instance.SubnetId),
	}

	return self
}

// mapInstanceToNodeName maps a EC2 instance to a k8s NodeName, by extracting the PrivateDNSName
func mapInstanceToNodeName(i *ec2.Instance) types.NodeName {
	return types.NodeName(_aws.StringValue(i.PrivateDnsName))
}
