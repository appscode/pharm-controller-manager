package softlayer

import (
	"context"
	"strings"

	"github.com/pharmer/cloud-controller-manager/cloud"
	"github.com/softlayer/softlayer-go/services"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

type loadbalancers struct {
	virtualServiceClient services.Virtual_Guest
	accountServiceClient services.Account
}

// newLoadbalancers returns a cloudprovider.LoadBalancer whose concrete type is a *loadbalancer.
func newLoadbalancers(virtualServiceClient services.Virtual_Guest,
	accountServiceClient services.Account) cloudprovider.LoadBalancer {
	return &loadbalancers{virtualServiceClient: virtualServiceClient,
		accountServiceClient: accountServiceClient}
}

// GetLoadBalancerName returns the name of the load balancer. Implementations must treat the
// *v1.Service parameter as read-only and not modify it.
func (l *loadbalancers) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	//GCE requires that the name of a load balancer starts with a lower case letter.
	ret := "a" + string(service.UID)
	ret = strings.Replace(ret, "-", "", -1)
	//AWS requires that the name of a load balancer is shorter than 32 bytes.
	if len(ret) > 32 {
		ret = ret[:32]
	}
	return ret
}

// GetLoadBalancer returns the *v1.LoadBalancerStatus of service.
//
// GetLoadBalancer will not modify service.
func (l *loadbalancers) GetLoadBalancer(_ context.Context, clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {
	return nil, false, cloud.ErrNotImplemented
}

// EnsureLoadBalancer ensures that the cluster is running a load balancer for
// service.
//
// EnsureLoadBalancer will not modify service or nodes.
func (l *loadbalancers) EnsureLoadBalancer(_ context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	return nil, cloud.ErrLBUnsupported

}

// UpdateLoadBalancer updates the load balancer for service to balance across
// the droplets in nodes.
//
// UpdateLoadBalancer will not modify service or nodes.
func (l *loadbalancers) UpdateLoadBalancer(_ context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return cloud.ErrLBUnsupported
}

// EnsureLoadBalancerDeleted deletes the specified loadbalancer if it exists.
// nil is returned if the load balancer for service does not exist or is
// successfully deleted.
//
// EnsureLoadBalancerDeleted will not modify service.
func (l *loadbalancers) EnsureLoadBalancerDeleted(_ context.Context, clusterName string, service *v1.Service) error {
	return cloud.ErrLBUnsupported
}
