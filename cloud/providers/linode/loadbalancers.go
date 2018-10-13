package linode

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/taoh/linodego"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
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
	annLinodeCheckPath       = "service.beta.kubernetes.io/linode-loadbalancer-check-path"
	annLinodeCheckBody       = "service.beta.kubernetes.io/linode-loadbalancer-check-body"
	annLinodeHealthCheckType = "service.beta.kubernetes.io/linode-loadbalancer-check-type"

	// annLinodeCertificateID is the annotation specifying the certificate ID
	// used for https protocol. This annoataion is required if annLinodeTLSPorts
	// is passed.
	annLinodeSSLCertificate = "service.beta.kubernetes.io/linode-loadbalancer-ssl-cert"
	annLinodeSSLKey         = "service.beta.kubernetes.io/linode-loadbalancer-ssl-key"

	annLinodeHealthCheckInterval = "service.beta.kubernetes.io/linode-loadbalancer-check-interval"
	annLinodeHealthCheckTimeout  = "service.beta.kubernetes.io/linode-loadbalancer-check-timeout"
	annLinodeHealthCheckAttempts = "service.beta.kubernetes.io/linode-loadbalancer-check-attempts"
	annLinodeHealthCheckPassive  = "service.beta.kubernetes.io/linode-loadbalancer-check-passive"

	annLinodeSessionPersistence = "service.beta.kubernetes.io/linode-loadbalancer-stickiness"

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
	zone   string
}

// newLoadbalancers returns a cloudprovider.LoadBalancer whose concrete type is a *loadbalancer.
func newLoadbalancers(client *linodego.Client, zone string) cloudprovider.LoadBalancer {
	return &loadbalancers{client: client, zone: zone}
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
func (l *loadbalancers) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {
	lbName := l.GetLoadBalancerName(ctx, clusterName, service)
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
}

// EnsureLoadBalancer ensures that the cluster is running a load balancer for
// service.
//
// EnsureLoadBalancer will not modify service or nodes.
func (l *loadbalancers) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	lbStatus, exists, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	if !exists {
		ip, err := l.buildLoadBalancerRequest(ctx, clusterName, service, nodes)
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

	err = l.UpdateLoadBalancer(ctx, clusterName, service, nodes)
	if err != nil {
		return nil, err
	}

	lbStatus, exists, err = l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	return lbStatus, nil

}

// UpdateLoadBalancer updates the load balancer for service to balance across
// the droplets in nodes.
//
// UpdateLoadBalancer will not modify service or nodes.
func (l *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	lbName := l.GetLoadBalancerName(ctx, clusterName, service)
	lb, err := l.lbByName(l.client, lbName)
	if err != nil {
		return err
	}

	nbs, err := l.client.NodeBalancerConfig.List(lb.NodeBalancerId, 0)
	if err != nil {
		return err
	}

	kubeNode := map[string]*v1.Node{}
	for _, node := range nodes {
		kubeNode[node.Name] = node
	}

	nodePort := map[int]v1.ServicePort{}
	for _, port := range service.Spec.Ports {
		nodePort[int(port.Port)] = port
	}

	for _, port := range service.Spec.Ports {
		found := false
		for _, nb := range nbs.NodeBalancerConfigs {
			if _, found := nodePort[nb.Port]; !found {
				if _, err = l.client.NodeBalancerConfig.Delete(lb.NodeBalancerId, nb.ConfigId); err != nil {
					return err
				}
				continue
			}
			if nb.Port == int(port.Port) {
				found = true
				protocol, err := getProtocol(service)
				if err != nil {
					return err
				}

				args := map[string]string{}
				args["Protocol"] = protocol
				args["Algorithm"] = getAlgorithm(service)
				healthArgs, err := getHealthCheck(service)
				if err != nil {
					return err
				}
				args = mergeMaps(args, healthArgs)
				tlsArgs, err := getTLSArgs(service, nb.Port, protocol)
				if err != nil {
					return err
				}
				args = mergeMaps(args, tlsArgs)
				_, err = l.client.NodeBalancerConfig.Update(nb.ConfigId, args)
				if err != nil {
					return err
				}

				nodeList, err := l.client.Node.List(nb.ConfigId, 0)
				for _, node := range nodeList.NodeBalancerNodes {
					if _, found := kubeNode[node.Label.String()]; !found {
						if _, err = l.client.Node.Delete(node.NodeId); err != nil {
							return err
						}
						continue
					}
					args := map[string]string{}
					args["Address"] = fmt.Sprintf("%v:%v", getNodeInternalIp(kubeNode[node.Label.String()]), port.NodePort)
					_, err := l.client.Node.Update(node.NodeId, args)
					if err != nil {
						return err
					}
				}
			}
		}
		if !found {
			ncid, err := l.createNodeBalancerConfig(service, lb.NodeBalancerId, int(port.Port))
			if err != nil {
				return err
			}
			for _, node := range nodes {
				if err = createNBNode(l.client, ncid, node, port.NodePort); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func createNBNode(client *linodego.Client, configID int, node *v1.Node, port int32) error {
	ip := fmt.Sprintf("%v:%v", getNodeInternalIp(node), port)
	args := map[string]string{}
	args["Weight"] = "100"
	args["Mode"] = "accept"

	_, err := client.Node.Create(configID, node.Name, ip, args)
	if err != nil {
		return err
	}
	return nil
}

// EnsureLoadBalancerDeleted deletes the specified loadbalancer if it exists.
// nil is returned if the load balancer for service does not exist or is
// successfully deleted.
//
// EnsureLoadBalancerDeleted will not modify service.
func (l *loadbalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	_, exists, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}
	lbName := l.GetLoadBalancerName(ctx, clusterName, service)
	lb, err := l.lbByName(l.client, lbName)
	if err != nil {
		return err
	}
	_, err = l.client.NodeBalancer.Delete(lb.NodeBalancerId)

	return err
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

func (l *loadbalancers) createNoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (int, error) {
	did, err := strconv.Atoi(l.zone)
	if err != nil {
		return -1, err
	}
	lbName := l.GetLoadBalancerName(ctx, clusterName, service)

	resp, err := l.client.NodeBalancer.Create(did, lbName, nil)
	if err != nil {
		return -1, err
	}
	return resp.NodeBalancerId.NodeBalancerId, nil
}

func (l *loadbalancers) createNodeBalancerConfig(service *v1.Service, nbId, port int) (int, error) {
	protocol, err := getProtocol(service)
	if err != nil {
		return -1, err
	}

	args := map[string]string{}
	args["Port"] = strconv.Itoa(port)
	args["Protocol"] = protocol
	args["Algorithm"] = getAlgorithm(service)
	if cp, ok := service.Annotations[annLinodeSessionPersistence]; ok {
		args["Stickiness"] = cp
	} else {
		args["Stickiness"] = "table"
	}

	healthArgs, err := getHealthCheck(service)
	if err != nil {
		return -1, err
	}
	args = mergeMaps(args, healthArgs)

	tlsArgs, err := getTLSArgs(service, port, protocol)
	if err != nil {
		return -1, err
	}
	args = mergeMaps(args, tlsArgs)
	resp, err := l.client.NodeBalancerConfig.Create(nbId, args)
	if err != nil {
		return -1, err
	}
	return resp.NodeBalancerConfigId.NodeBalancerConfigId, nil
}

func getHealthCheck(service *v1.Service) (map[string]string, error) {
	args := map[string]string{}

	health, err := getHealthCheckType(service)
	if err != nil {
		return args, nil
	}
	args["check"] = health
	if health == "http" || health == "http_body" {
		path := service.Annotations[annLinodeCheckPath]
		if path == "" {
			path = "/"
		}
		args["check_path"] = path
	}

	if health == "http_body" {
		body := service.Annotations[annLinodeCheckBody]
		if body == "" {
			return args, fmt.Errorf("for health check type http_body need body regex annotation %v", annLinodeCheckBody)
		}
		args["check_body"] = body
	}
	if ci, ok := service.Annotations[annLinodeHealthCheckInterval]; ok {
		args["check_interval"] = ci
	} else {
		args["check_interval"] = "5"
	}

	if ct, ok := service.Annotations[annLinodeHealthCheckTimeout]; ok {
		args["check_timeout"] = ct
	} else {
		args["check_timeout"] = "3"
	}

	if ca, ok := service.Annotations[annLinodeHealthCheckAttempts]; ok {
		args["check_attempts"] = ca
	} else {
		args["check_attempts"] = "2"
	}

	if cp, ok := service.Annotations[annLinodeHealthCheckPassive]; ok {
		args["check_passive"] = cp
	} else {
		args["check_passive"] = "true"
	}

	return args, nil
}

func getTLSArgs(service *v1.Service, port int, protocol string) (map[string]string, error) {
	args := map[string]string{}
	tlsPorts, err := getTLSPorts(service)
	if err != nil {
		return args, err
	}
	if len(tlsPorts) > 0 {
		for _, tlsPort := range tlsPorts {
			if tlsPort == port && protocol == "https" {
				cert, key := getSSLCertInfo(service)
				if cert == "" && key == "" {
					return args, fmt.Errorf("must set %v and %v annotation for https protocol", annLinodeSSLCertificate, annLinodeSSLKey)
				}
				if cert != "" {
					args["ssl_cert"] = cert
				}
				if key != "" {
					args["ssl_key"] = key
				}

			}
		}
	}

	return args, nil
}

func mergeMaps(first, second map[string]string) map[string]string {
	for k, v := range second {
		first[k] = v
	}
	return first
}

// buildLoadBalancerRequest returns a *godo.LoadBalancerRequest to balance
// requests for service across nodes.
func (l *loadbalancers) buildLoadBalancerRequest(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (string, error) {
	lb, err := l.createNoadBalancer(ctx, clusterName, service)
	if err != nil {
		return "", err
	}

	nb, err := l.client.NodeBalancer.List(lb)
	if err != nil {
		return "", err
	}
	if len(nb.NodeBalancer) == 0 {
		return "", fmt.Errorf("nodebalancer with id %v not found", lb)
	}
	ports := service.Spec.Ports
	for _, port := range ports {

		ncid, err := l.createNodeBalancerConfig(service, lb, int(port.Port))
		if err != nil {
			return "", err
		}

		for _, node := range nodes {
			if err = createNBNode(l.client, ncid, node, port.NodePort); err != nil {
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

func getHealthCheckType(service *v1.Service) (string, error) {
	hType, ok := service.Annotations[annLinodeHealthCheckType]
	if !ok {
		return "connection", nil
	}
	if hType != "connection" && hType != "http" && hType != "http_body" {
		return "", fmt.Errorf("invalid health check type: %q specifed in annotation: %q", hType, annLinodeHealthCheckType)
	}
	return hType, nil
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

func getNodeInternalIp(node *v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}

// getAlgorithm returns the load balancing algorithm to use for service.
// round_robin is returned when service does not specify an algorithm.
func getAlgorithm(service *v1.Service) string {
	algo := service.Annotations[annLinodeAlgorithm]

	switch algo {
	case "least_connections":
		return "leastconn"
	case "source":
		return "source"
	default:
		return "roundrobin"
	}
}

func getSSLCertInfo(service *v1.Service) (string, string) {
	cert := service.Annotations[annLinodeSSLCertificate]
	if cert != "" {
		cb, _ := base64.StdEncoding.DecodeString(cert)
		cert = string(cb)
	}
	key := service.Annotations[annLinodeSSLKey]
	if key != "" {
		kb, _ := base64.StdEncoding.DecodeString(key)
		key = string(kb)
	}
	return cert, key
}
