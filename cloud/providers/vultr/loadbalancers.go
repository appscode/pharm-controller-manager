package vultr

import (
	"errors"

	gv "github.com/JamesClonk/vultr/lib"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

var lbNotFound = errors.New("loadbalancer not found")

type loadbalancers struct {
	client *gv.Client
}

// newLoadbalancers returns a cloudprovider.LoadBalancer whose concrete type is a *loadbalancer.
func newLoadbalancers(client *gv.Client) cloudprovider.LoadBalancer {
	return &loadbalancers{client: client}
}

// GetLoadBalancer returns the *v1.LoadBalancerStatus of service.
//
// GetLoadBalancer will not modify service.
func (l *loadbalancers) GetLoadBalancer(clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {
	return nil, false, lbNotFound
}

// EnsureLoadBalancer ensures that the cluster is running a load balancer for
// service.
//
// EnsureLoadBalancer will not modify service or nodes.
func (l *loadbalancers) EnsureLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	return nil, lbNotFound

}

// UpdateLoadBalancer updates the load balancer for service to balance across
// the droplets in nodes.
//
// UpdateLoadBalancer will not modify service or nodes.
func (l *loadbalancers) UpdateLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return lbNotFound
}

// EnsureLoadBalancerDeleted deletes the specified loadbalancer if it exists.
// nil is returned if the load balancer for service does not exist or is
// successfully deleted.
//
// EnsureLoadBalancerDeleted will not modify service.
func (l *loadbalancers) EnsureLoadBalancerDeleted(clusterName string, service *v1.Service) error {
	return lbNotFound
}
