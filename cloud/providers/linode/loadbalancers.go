package linode

import (
	"github.com/pharmer/cloud-controller-manager/cloud"
	"github.com/taoh/linodego"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"errors"
	"strconv"
	"fmt"
	"strings"
)


const (
	// annDOProtocol is the annotation used to specify the default protocol
	// for DO load balancers. For ports specifed in annDOTLSPorts, this protocol
	// is overwritten to https. Options are tcp, http and https. Defaults to tcp.
	annLinodeProtocol = "service.beta.kubernetes.io/linode-loadbalancer-protocol"

	// annLinodeTLSPorts is the annotation used to specify which ports of the loadbalancer
	// should use the https protocol. This is a comma separated list of ports
	// (e.g. 443,6443,7443).
	annLinodeTLSPorts = "service.beta.kubernetes.io/linode-loadbalancer-tls-ports"

	// annLinodeTLSPassThrough is the annotation used to specify whether the
	// Linode loadbalancer should pass encrypted data to backend droplets.
	// This is optional and defaults to false.
	annLinodeTLSPassThrough = "service.beta.kubernetes.io/linode-loadbalancer-tls-passthrough"

	// annLinodeCertificateID is the annotation specifying the certificate ID
	// used for https protocol. This annoataion is required if annLinodeTLSPorts
	// is passed.
	annLinodeCertificateID = "service.beta.kubernetes.io/linode-loadbalancer-certificate-id"

	// annLinodeAlgorithm is the annotation specifying which algorithm Linode loadbalancer
	// should use. Options are round_robin and least_connections. Defaults
	// to round_robin.
	annLinodeAlgorithm = "service.beta.kubernetes.io/linode-loadbalancer-algorithm"

	// defaultActiveTimeout is the number of seconds to wait for a load balancer to
	// reach the active state.
	defaultActiveTimeout = 90

	// defaultActiveCheckTick is the number of seconds between load balancer
	// status checks when waiting for activation.
	defaultActiveCheckTick = 5

	// statuses for Digital Ocean load balancer
	lbStatusNew     = "new"
	lbStatusActive  = "active"
	lbStatusErrored = "errored"
)
var lbNotFound = errors.New("loadbalancer not found")


type loadbalancers struct {
	client *linodego.Client
	zone string
}

// newLoadbalancers returns a cloudprovider.LoadBalancer whose concrete type is a *loadbalancer.
func newLoadbalancers(client *linodego.Client, zone string) cloudprovider.LoadBalancer {
	return &loadbalancers{client: client, zone:zone}
}

// GetLoadBalancer returns the *v1.LoadBalancerStatus of service.
//
// GetLoadBalancer will not modify service.
func (l *loadbalancers) GetLoadBalancer(clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {
	lbName := cloudprovider.GetLoadBalancerName(service)
	lb, err := l.lbByName(l.client, lbName)
	if err != nil {
		if err == lbNotFound {
			return nil, false, nil
		}

		return nil, false, err
	}
	/*if lb.Status != lbStatusActive {
		lb, err = l.waitActive(lb.ID)
		if err != nil {
			return nil, true, fmt.Errorf("error waiting for load balancer to be active %v", err)
		}
	}*/

	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: lb.Address4,
			},
		},
	}, true, nil
	return nil, false, cloud.ErrNotImplemented
}

// EnsureLoadBalancer ensures that the cluster is running a load balancer for
// service.
//
// EnsureLoadBalancer will not modify service or nodes.
func (l *loadbalancers) EnsureLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	lbStatus, exists, err := l.GetLoadBalancer(clusterName, service)
	if err != nil {
		return nil, err
	}

	if !exists {
		ip,  err := l.buildLoadBalancerRequest(service, nodes)
		if err != nil {
			return nil, err
		}

		return &v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{
				{
					IP: ip,
				},
			},
		}, nil
	}

	err = l.UpdateLoadBalancer(clusterName, service, nodes)
	if err != nil {
		return nil, err
	}

	lbStatus, exists, err = l.GetLoadBalancer(clusterName, service)
	if err != nil {
		return nil, err
	}

	return lbStatus, nil

}

// UpdateLoadBalancer updates the load balancer for service to balance across
// the droplets in nodes.
//
// UpdateLoadBalancer will not modify service or nodes.
func (l *loadbalancers) UpdateLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return cloud.ErrLBUnsupported
}

// EnsureLoadBalancerDeleted deletes the specified loadbalancer if it exists.
// nil is returned if the load balancer for service does not exist or is
// successfully deleted.
//
// EnsureLoadBalancerDeleted will not modify service.
func (l *loadbalancers) EnsureLoadBalancerDeleted(clusterName string, service *v1.Service) error {
	return cloud.ErrLBUnsupported
}


// lbByName gets a DigitalOcean Load Balancer by name. The returned error will
// be lbNotFound if the load balancer does not exist.
func (l *loadbalancers) lbByName(client *linodego.Client, name string) (*linodego.LinodeNodeBalancer, error) {
	lbs, err := l.client.NodeBalancer.List(0)
	if err != nil {
		return nil, err
	}

	for _, lb := range lbs.NodeBalancer {
		if lb.Label.String() == name {
			return &lb, nil
		}
	}

	return nil, lbNotFound
}

func (l *loadbalancers) createNoadBalancer(service *v1.Service) (int, error) {
	did, err := strconv.Atoi(l.zone)
	if err != nil {
		return -1, err
	}
	lbName := cloudprovider.GetLoadBalancerName(service)

	resp, err := l.client.NodeBalancer.Create(did, lbName, nil)
	if err != nil {
		return -1, err
	}
	return resp.NodeBalancerId.NodeBalancerId, nil
}

func (l *loadbalancers) createNodeBalancerConfig(nbId int, args map[string]string) (int, error)  {
	resp, err := l.client.NodeBalancerConfig.Create(nbId, args)
	if err != nil {
		return -1, err
	}
	return resp.NodeBalancerConfigId.NodeBalancerConfigId, nil
}


// buildLoadBalancerRequest returns a *godo.LoadBalancerRequest to balance
// requests for service across nodes.
func (l *loadbalancers) buildLoadBalancerRequest(service *v1.Service, nodes []*v1.Node)  (string, error){
	lb, err := l.createNoadBalancer(service)
	if err != nil {
		return  "", err
	}

	nb, err := l.client.NodeBalancer.List(lb)
	if err != nil {
		return "", err
	}
	ports := service.Spec.Ports
	for _, port := range ports {
		protocol, err := getProtocol(service)
		if err != nil {
			return  "", err
		}

		args := map[string]string{}
		args["Port"] = strconv.Itoa(int(port.Port))
		args["Protocol"] = protocol
		args["Algorithm"] = "roundrobin"
		args["Stickiness"] = "table"
		args["check"] = "connection"
		args["check_interval"] = "5"
		args["check_timeout"] = "3"
		args["check_attempts"] = "2"
		args["check_passive"] = "true"
		ncid, err := l.createNodeBalancerConfig(lb, args)
		if err != nil {
			return  "", err
		}


		for _, node := range nodes {
			_, err := l.client.Node.Create(ncid, node.Name, fmt.Sprint("%v:%v", nb.NodeBalancer[0].Address4, port.NodePort), nil )
			if err != nil {
				return "", err
			}
		}
	}
	return nb.NodeBalancer[0].Address4, nil
}


// getProtocol returns the desired protocol of service.
func getProtocol(service *v1.Service) (string, error) {
	protocol, ok := service.Annotations[annLinodeProtocol]
	if !ok {
		return "tcp", nil
	}

	if protocol != "tcp" && protocol != "http" && protocol != "https" {
		return "", fmt.Errorf("invalid protocol: %q specifed in annotation: %q", protocol, annLinodeProtocol)
	}

	return protocol, nil
}

// getTLSPorts returns the ports of service that are set to use TLS.
func getTLSPorts(service *v1.Service) ([]int, error) {
	tlsPorts, ok := service.Annotations[annLinodeTLSPorts]
	if !ok {
		return nil, nil
	}

	tlsPortsSlice := strings.Split(tlsPorts, ",")

	tlsPortsInt := make([]int, len(tlsPortsSlice))
	for i, tlsPort := range tlsPortsSlice {
		port, err := strconv.Atoi(tlsPort)
		if err != nil {
			return nil, err
		}

		tlsPortsInt[i] = port
	}

	return tlsPortsInt, nil
}