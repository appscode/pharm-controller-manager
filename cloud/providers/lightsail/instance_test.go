package lightsail

import (
	"fmt"
	"testing"

	_aws "github.com/aws/aws-sdk-go/aws"

	//"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

func TestInstances(t *testing.T) {
	client := getClient()
	token := ""
	instances, err := client.GetInstances(&lightsail.GetInstancesInput{PageToken: &token})
	fmt.Println(err)
	fmt.Println(instances, instances.NextPageToken)
}

func TestInstance(t *testing.T) {
	client := getClient()
	instanceName := "ls5-master"
	ins, err := client.GetInstance(&lightsail.GetInstanceInput{
		InstanceName: _aws.String(instanceName),
	})
	fmt.Println(err)
	fmt.Println(ins.Instance)
}

func getClient() *lightsail.Lightsail {
	region := "us-west-2"
	conf := &_aws.Config{
		Region: &region,
		//Credentials: credentials.NewStaticCredentials("", "", ""),
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}
	return lightsail.New(sess)
}
