package lightsail

import "github.com/aws/aws-sdk-go/service/lightsail"

func allInstanceList(client *lightsail.Lightsail) ([]*lightsail.Instance, error) {
	list := []*lightsail.Instance{}

	nextPageToken := ""
	for {
		instances, err := client.GetInstances(&lightsail.GetInstancesInput{
			PageToken: &nextPageToken,
		})
		if err != nil {
			return nil, err
		}
		list = append(list, instances.Instances...)
		if instances.NextPageToken == nil {
			break
		}
		nextPageToken = *instances.NextPageToken
	}
	return list, nil
}
