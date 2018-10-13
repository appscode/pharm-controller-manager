package linode

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/taoh/linodego"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

var _ cloudprovider.LoadBalancer = new(loadbalancers)

type fakeNodeBalancerService struct {
	createFn func(int, string, map[string]string) (*linodego.NodeBalancerResponse, error)
	deleteFn func(int) (*linodego.NodeBalancerResponse, error)
	listFn   func(int) (*linodego.NodeBalancerListResponse, error)
	updateFn func(int, map[string]string) (*linodego.NodeBalancerResponse, error)
}

func (f *fakeNodeBalancerService) Create(datacenterId int, label string, args map[string]string) (*linodego.NodeBalancerResponse, error) {
	return f.createFn(datacenterId, label, args)
}

func (f *fakeNodeBalancerService) Delete(nodeBalancerId int) (*linodego.NodeBalancerResponse, error) {
	return f.deleteFn(nodeBalancerId)
}

func (f *fakeNodeBalancerService) List(nodeBalancerId int) (*linodego.NodeBalancerListResponse, error) {
	return f.listFn(nodeBalancerId)
}

func (f *fakeNodeBalancerService) Update(nodeBalancerId int, args map[string]string) (*linodego.NodeBalancerResponse, error) {
	return f.updateFn(nodeBalancerId, args)
}

type fakeNodeBalancerConfigService struct {
	createFn func(int, map[string]string) (*linodego.NodeBalancerConfigResponse, error)
	deleteFn func(int, int) (*linodego.NodeBalancerConfigResponse, error)
	listFn   func(int, int) (*linodego.NodeBalancerConfigListResponse, error)
	updateFn func(int, map[string]string) (*linodego.NodeBalancerConfigResponse, error)
}

func (f *fakeNodeBalancerConfigService) Create(nodeBalancerId int, args map[string]string) (*linodego.NodeBalancerConfigResponse, error) {
	return f.createFn(nodeBalancerId, args)
}

func (f *fakeNodeBalancerConfigService) List(nodeBalancerId int, configId int) (*linodego.NodeBalancerConfigListResponse, error) {
	return f.listFn(nodeBalancerId, configId)
}

func (f *fakeNodeBalancerConfigService) Delete(nodeBalancerId int, configId int) (*linodego.NodeBalancerConfigResponse, error) {
	return f.deleteFn(nodeBalancerId, configId)
}

func (f *fakeNodeBalancerConfigService) Update(configId int, args map[string]string) (*linodego.NodeBalancerConfigResponse, error) {
	return f.updateFn(configId, args)
}

type fakeNodeBalancerNodeService struct {
	createFn func(int, string, string, map[string]string) (*linodego.NodeBalancerNodeResponse, error)
	deleteFn func(int) (*linodego.NodeBalancerNodeResponse, error)
	listFn   func(int, int) (*linodego.NodeBalancerNodeListResponse, error)
	updateFn func(int, map[string]string) (*linodego.NodeBalancerNodeResponse, error)
}

func (f *fakeNodeBalancerNodeService) Create(configId int, label string, address string, args map[string]string) (*linodego.NodeBalancerNodeResponse, error) {
	return f.createFn(configId, label, address, args)
}

func (f *fakeNodeBalancerNodeService) Delete(nodeId int) (*linodego.NodeBalancerNodeResponse, error) {
	return f.deleteFn(nodeId)
}

func (f *fakeNodeBalancerNodeService) List(configId int, nodeId int) (*linodego.NodeBalancerNodeListResponse, error) {
	return f.listFn(configId, nodeId)
}

func (f *fakeNodeBalancerNodeService) Update(NodeId int, args map[string]string) (*linodego.NodeBalancerNodeResponse, error) {
	return f.updateFn(NodeId, args)
}
func newFakeLBClient(fakeNB *fakeNodeBalancerService, fakeLinode *fakeLinodeService, fakeNBC *fakeNodeBalancerConfigService, fakeN *fakeNodeBalancerNodeService) *linodego.Client {
	client := linodego.NewClient("", nil)
	client.Linode = fakeLinode
	client.NodeBalancer = fakeNB
	client.NodeBalancerConfig = fakeNBC
	client.Node = fakeN
	return client
}

func Test_getAlgorithm(t *testing.T) {
	testcases := []struct {
		name      string
		service   *v1.Service
		algorithm string
	}{
		{
			"algorithm should be least_connection",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeAlgorithm: "least_connections",
					},
				},
			},
			"leastconn",
		},
		{
			"algorithm should be source",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeAlgorithm: "source",
					},
				},
			},
			"source",
		},
		{
			"algorithm should be round_robin",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeAlgorithm: "roundrobin",
					},
				},
			},
			"roundrobin",
		},
		{
			"invalid algorithm should default to round_robin",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeAlgorithm: "invalid",
					},
				},
			},
			"roundrobin",
		},
		{
			"no algorithm specified should default to round_robin",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
				},
			},
			"roundrobin",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			algorithm := getAlgorithm(test.service)
			if algorithm != test.algorithm {
				t.Error("unexpected algoritmh")
				t.Logf("expected: %q", test.algorithm)
				t.Logf("actual: %q", algorithm)
			}
		})
	}
}

func Test_getCertificate(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIICuDCCAaCgAwIBAgIBADANBgkqhkiG9w0BAQsFADANMQswCQYDVQQDEwJjYTAe
Fw0xNzExMjcwNTQ3NDJaFw0yNzExMjUwNTQ3NDJaMA0xCzAJBgNVBAMTAmNhMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7GdQywI4pm50c0TyiOoKi4ar
AwSSgHdDSQFNM4k2ssXuem8S1DMRScY663LYn14n1PM6fppCtZWC/vtsDnmEEGUy
/w+hJ8w90uFExMBmkn8D765W59jWtE3x3/7Kd0PGyiXGsdqRxmhainOO6p9Q8/Ln
SwPpsVMRnbSDAnoNqRFK59YIfxoQXML2+e45M+oFbxUoi2xXQCsj1qdxTshtqwT/
7u0nWOOSoq8a3YKv7zk+qZwCNe0PSKXKbnNNJgzdx+UJWBChvrt0Ndm+swTG125B
lMlBrmNJOYWdNGLKuFsWX+OPC7fNj9VwxarOy+H5ykLH0i+7jxCpgYGF+eFDvwID
AQABoyMwITAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG
9w0BAQsFAAOCAQEAJwH7LC0d1z7r/ztQ2AekAsqwu+So/EHqVzGRP4jHRJMPtYN1
SFBCIdLYmACuj0MfWLyPy2WkpT4F3h9CrQzswl42mdjd0Q5vutjGsDo6Jta9Jf4Y
ouM2felPMvbAXHIAFIXXa64okcnWoJzp1oAxfCieqZXb2yhPJMcULtNUC5NtYEpG
oNF1FzyoGh5GNpeARDnzU7RACF9PiCxx8hWHV9V09IXXP5TjBDdc4rvll7P93W7V
3WV87/Aeh/W8TueGYBeUOmzn63VbEkpmGT9KJe8t+IrVymuG4rYS08z6g5Ib9FNh
KHB9fdnWTibkrKB/319X4GfMjGNN2/YyER2F8g==
-----END CERTIFICATE-----`
	key := `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA7GdQywI4pm50c0TyiOoKi4arAwSSgHdDSQFNM4k2ssXuem8S
1DMRScY663LYn14n1PM6fppCtZWC/vtsDnmEEGUy/w+hJ8w90uFExMBmkn8D765W
59jWtE3x3/7Kd0PGyiXGsdqRxmhainOO6p9Q8/LnSwPpsVMRnbSDAnoNqRFK59YI
fxoQXML2+e45M+oFbxUoi2xXQCsj1qdxTshtqwT/7u0nWOOSoq8a3YKv7zk+qZwC
Ne0PSKXKbnNNJgzdx+UJWBChvrt0Ndm+swTG125BlMlBrmNJOYWdNGLKuFsWX+OP
C7fNj9VwxarOy+H5ykLH0i+7jxCpgYGF+eFDvwIDAQABAoIBAQDj8sdDyPOI/66H
y261uD7MxOC2+zysZNNbXMbtL5yviw1lvx5/wHImGd+MUmQwX2C3BIVduC8k2nLC
nPpXhrJiAMLIkHCLaHQgmBhwQzlkftbz0L55tmto1lOo8gyWLaNMHlrV+fRgRRUw
tTaUY2RypcCCY9Z9pqSw1XMR+1CauHhicfY9K1rQgF8xtZ6sB+P7y2SwVlp2OjBr
R7E66O4s3LPf6A30ZbnaertZrrO36//sXKMKLeUlginzE3oMZBfr1IMYtd+5JKVX
axyMMNAqUjdpJk/ahE0B52Toebj9XSxTNkiswmNS6Zve9CV5oiRkntsDZXpiDnRb
7lEHXnjhAoGBAPtYQ+Y+sg4utk4BOIK2apjUVLwXuDQiCREzCnhA3CLCSqJMb6Y8
7N1+KzRZYeDNECt5DOJOrUqM2pTIQ+RkZEhaUfJr1ILFGQmD7FhjxrM8nQh5gUKO
9fGEKPPIOshkUoVCNm5HMixa7YnGM1xhvXvHLPSXILwuz082e2ZnI5SRAoGBAPDI
NSWEJ3d81YnIK6aDoPmpDv0FG+TweYqIdEs8eja8TN7Bpbx2vuUS/vkWsjJeyTkS
7V0Bq6bKVwfiFCYjEPNQ8qekifb+tHRLu6DRbj4UbeAcZXr3C5mcUQk07/84gXXj
FUDfT8EI6Eerr6RM75CTN7nesiwGXMjyYSSomTtPAoGBAOs8s+fVO95sN7GQEOy9
f8zjxR55cKxSQnw3chAUXDOn9iQqN8C1etbeU99d3G6CXiTh2X4hNqz0YUsol+o1
T2osJlAmPbHaeFFgiB492+U60Jny5lh95o+RKqbm+qU8x8LysnDJ75p1y6XLu5w1
2hrz0g5lN30IrnwruJih5ToRAoGBAISK8RaRxNf1k+aglca3tqk38tQ9N7my1nT3
4Gx6AhyXUwlcN8ui4jpfVpPvdnBb1RDh5l/IR6EsyPPB8616qB4IdUrrPDcGxncu
KT7BipoJzOINP5+M1oncjo8u4N3xUPJ/6ncndlOgf5zUWX9sCoPfRlG+0P2DExha
tDblyFPpAoGAC29vNqFODcEhmiXfOEW2pvMXztgBVKopYGnKl2qDv9uBtMkn7qUi
V/teWB7SHT+secFM8H2lukmIO6PQhqH3sK7mAGxGXWWBUMVKeU8assuJmqXQsvMs
b8QPmGZdja1VyGqpAMkPmQOu9N5RbhKw1UOU/XGa31p6v96oayL+u8Q=
-----END RSA PRIVATE KEY-----`
	testcases := []struct {
		name    string
		service *v1.Service
		cert    string
		key     string
	}{
		{
			"certificate set",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeSSLCertificate: base64.StdEncoding.EncodeToString([]byte(cert)),
						annLinodeSSLKey:         base64.StdEncoding.EncodeToString([]byte(key)),
					},
				},
			},
			cert,
			key,
		},
		{
			"certificate not set",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					UID:         "abc123",
					Annotations: map[string]string{},
				},
			},
			"",
			"",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			c, k := getSSLCertInfo(test.service)
			if c != test.cert {
				t.Error("unexpected certificate")
				t.Logf("expected: %q", test.cert)
				t.Logf("actual: %q", c)
			}
			if k != test.key {
				t.Error("unexpected key")
				t.Logf("expected: %q", test.key)
				t.Logf("actual: %q", k)
			}
		})
	}
}

func Test_getTLSPorts(t *testing.T) {
	testcases := []struct {
		name     string
		service  *v1.Service
		tlsPorts []int
		err      error
	}{
		{
			"tls port specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeTLSPorts: "443",
					},
				},
			},
			[]int{443},
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			tlsPorts, err := getTLSPorts(test.service)
			if !reflect.DeepEqual(tlsPorts, test.tlsPorts) {
				t.Error("unexpected TLS ports")
				t.Logf("expected %v", test.tlsPorts)
				t.Logf("actual: %v", tlsPorts)
			}

			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_getProtocol(t *testing.T) {
	testcases := []struct {
		name     string
		service  *v1.Service
		protocol string
		err      error
	}{
		{
			"no protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
				},
			},
			"tcp",
			nil,
		},
		{
			"tcp protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeProtocol: "http",
					},
				},
			},
			"http",
			nil,
		},
		{
			"invalid protocol",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeProtocol: "invalid",
					},
				},
			},
			"",
			fmt.Errorf("invalid protocol: %q specifed in annotation: %q", "invalid", annLinodeProtocol),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			protocol, err := getProtocol(test.service)
			if protocol != test.protocol {
				t.Error("unexpected protocol")
				t.Logf("expected: %q", test.protocol)
				t.Logf("actual: %q", protocol)
			}

			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %q", test.err)
				t.Logf("actual: %q", err)
			}
		})
	}
}

func Test_getHealthCheckType(t *testing.T) {
	testcases := []struct {
		name       string
		service    *v1.Service
		healthType string
		err        error
	}{
		{
			"no type specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					UID:         "abc123",
					Annotations: map[string]string{},
				},
			},
			"connection",
			nil,
		},
		{
			"http specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeHealthCheckType: "http",
					},
				},
			},
			"http",
			nil,
		},
		{
			"invalid specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annLinodeHealthCheckType: "invalid",
					},
				},
			},
			"",
			fmt.Errorf("invalid health check type: %q specifed in annotation: %q", "invalid", annLinodeHealthCheckType),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			hType, err := getHealthCheckType(test.service)
			if !reflect.DeepEqual(hType, test.healthType) {
				t.Error("unexpected health check type")
				t.Logf("expected: %v", test.healthType)
				t.Logf("actual: %v", hType)
			}

			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_getNodeInternalIp(t *testing.T) {
	testcases := []struct {
		name    string
		node    *v1.Node
		address string
	}{
		{
			"node internal ip specified",
			&v1.Node{
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							v1.NodeInternalIP,
							"127.0.0.1",
						},
					},
				},
			},
			"127.0.0.1",
		},
		{
			"node internal ip not specified",
			&v1.Node{
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							v1.NodeExternalIP,
							"127.0.0.1",
						},
					},
				},
			},
			"",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			ip := getNodeInternalIp(test.node)
			if ip != test.address {
				t.Error("unexpected certificate")
				t.Logf("expected: %q", test.address)
				t.Logf("actual: %q", ip)
			}
		})
	}

}

func Test_createNoadBalancer(t *testing.T) {
	testcases := []struct {
		name     string
		createFn func(int, string, map[string]string) (*linodego.NodeBalancerResponse, error)
		service  *v1.Service
		nbId     int
		err      error
	}{
		{
			"create nodebalancer",
			func(i int, s string, strings map[string]string) (*linodego.NodeBalancerResponse, error) {
				return &linodego.NodeBalancerResponse{
					newFakeOKResponse("nodebalancer.create"),
					linodego.LinodeNodeBalancerId{
						12345,
					},
				}, nil
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					UID:         "foobar123",
					Annotations: map[string]string{},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			12345,
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeNB := &fakeNodeBalancerService{
				createFn: test.createFn,
			}
			fakeClient := newFakeLBClient(fakeNB, nil, nil, nil)
			lb := &loadbalancers{fakeClient, "1"}
			id, err := lb.createNoadBalancer(context.Background(), "k1", test.service)
			if id != test.nbId {
				t.Error("unexpected nodeID")
				t.Logf("expected: %v", test.nbId)
				t.Logf("actual: %v", id)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_buildLoadBalancerRequest(t *testing.T) {
	testcases := []struct {
		name        string
		createNBFn  func(int, string, map[string]string) (*linodego.NodeBalancerResponse, error)
		listFn      func(int) (*linodego.NodeBalancerListResponse, error)
		createNBCFn func(int, map[string]string) (*linodego.NodeBalancerConfigResponse, error)
		createNFn   func(int, string, string, map[string]string) (*linodego.NodeBalancerNodeResponse, error)
		service     *v1.Service
		nodes       []*v1.Node
		address     string
		err         error
	}{
		{
			"build load balancer",
			func(i int, s string, strings map[string]string) (*linodego.NodeBalancerResponse, error) {
				return &linodego.NodeBalancerResponse{
					newFakeOKResponse("nodebalancer.create"),
					linodego.LinodeNodeBalancerId{
						12345,
					},
				}, nil
			},
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{
						{
							NodeBalancerId: i,
							Address4:       "127.0.0.1",
						},
					},
				}, nil
			},
			func(i int, strings map[string]string) (*linodego.NodeBalancerConfigResponse, error) {
				return &linodego.NodeBalancerConfigResponse{
					newFakeOKResponse("nodebalancer.config.create"),
					linodego.NodeBalancerConfigId{32104},
				}, nil
			},
			func(i int, s string, s2 string, strings map[string]string) (*linodego.NodeBalancerNodeResponse, error) {
				return &linodego.NodeBalancerNodeResponse{
					newFakeOKResponse("nodebalancer.node.create"),
					linodego.NodeBalancerNodeId{123},
				}, nil
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			[]*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-3",
					},
				},
			},
			"127.0.0.1",
			nil,
		},
		{
			"failed to build load balancer",
			func(i int, s string, strings map[string]string) (*linodego.NodeBalancerResponse, error) {
				return &linodego.NodeBalancerResponse{
					newFakeOKResponse("nodebalancer.create"),
					linodego.LinodeNodeBalancerId{
						12345,
					},
				}, nil
			},
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{},
				}, nil
			},
			func(i int, strings map[string]string) (*linodego.NodeBalancerConfigResponse, error) {
				return &linodego.NodeBalancerConfigResponse{
					newFakeOKResponse("nodebalancer.config.create"),
					linodego.NodeBalancerConfigId{32104},
				}, nil
			},
			func(i int, s string, s2 string, strings map[string]string) (*linodego.NodeBalancerNodeResponse, error) {
				return &linodego.NodeBalancerNodeResponse{
					newFakeOKResponse("nodebalancer.node.create"),
					linodego.NodeBalancerNodeId{123},
				}, nil
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			[]*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-3",
					},
				},
			},
			"",
			fmt.Errorf("nodebalancer with id %v not found", 12345),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeNB := &fakeNodeBalancerService{
				createFn: test.createNBFn,
				listFn:   test.listFn,
			}
			fakeNBC := &fakeNodeBalancerConfigService{
				createFn: test.createNBCFn,
			}
			fakeN := &fakeNodeBalancerNodeService{
				createFn: test.createNFn,
			}
			fakeClient := newFakeLBClient(fakeNB, nil, fakeNBC, fakeN)
			lb := &loadbalancers{fakeClient, "1"}
			id, err := lb.buildLoadBalancerRequest(context.Background(), "k1", test.service, test.nodes)
			if id != test.address {
				t.Error("unexpected nodeID")
				t.Logf("expected: %v", test.address)
				t.Logf("actual: %v", id)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_EnsureLoadBalancerDeleted(t *testing.T) {
	testcases := []struct {
		name        string
		deleteFn    func(int) (*linodego.NodeBalancerResponse, error)
		listFn      func(int) (*linodego.NodeBalancerListResponse, error)
		clusterName string
		service     *v1.Service
		err         error
	}{
		{
			"load balancer delete",
			func(i int) (*linodego.NodeBalancerResponse, error) {
				return &linodego.NodeBalancerResponse{
					newFakeOKResponse("nodebalancer.delete"),
					linodego.LinodeNodeBalancerId{i},
				}, nil
			},
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				label := linodego.CustomString{}
				err := label.UnmarshalJSON([]byte("afoobar123"))
				if err != nil {
					return nil, err
				}
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{
						{
							// loadbalancer names are a + service.UID
							// see cloudprovider.GetLoadBalancerName
							Label: label,
						},
					},
				}, nil
			},
			"linodelb",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			nil,
		},
		{
			"load balancer not exists",
			func(i int) (*linodego.NodeBalancerResponse, error) {
				return &linodego.NodeBalancerResponse{
					newFakeNotOKResponse(400, "node balancer not exists", "nodebalancer.delete"),
					linodego.LinodeNodeBalancerId{},
				}, fmt.Errorf("node balancer not exists")
			},
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				label := linodego.CustomString{}
				err := label.UnmarshalJSON([]byte("afoobar321"))
				if err != nil {
					return nil, err
				}
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{
						{
							// loadbalancer names are a + service.UID
							// see cloudprovider.GetLoadBalancerName
							Label: label,
						},
					},
				}, nil
			},
			"linodelb",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeNB := &fakeNodeBalancerService{
				deleteFn: test.deleteFn,
				listFn:   test.listFn,
			}

			fakeClient := newFakeLBClient(fakeNB, nil, nil, nil)
			lb := &loadbalancers{fakeClient, "1"}
			err := lb.EnsureLoadBalancerDeleted(context.Background(), test.clusterName, test.service)

			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_EnsureLoadBalancer(t *testing.T) {
	testcases := []struct {
		name       string
		createNBFn func(int, string, map[string]string) (*linodego.NodeBalancerResponse, error)
		listNBFn   func(int) (*linodego.NodeBalancerListResponse, error)

		createNBCFn func(int, map[string]string) (*linodego.NodeBalancerConfigResponse, error)
		listNBCFn   func(int, int) (*linodego.NodeBalancerConfigListResponse, error)
		deleteNBCFn func(int, int) (*linodego.NodeBalancerConfigResponse, error)
		updateNBCFn func(int, map[string]string) (*linodego.NodeBalancerConfigResponse, error)

		deleteNFn func(int) (*linodego.NodeBalancerNodeResponse, error)
		updateNFn func(int, map[string]string) (*linodego.NodeBalancerNodeResponse, error)
		createNFn func(int, string, string, map[string]string) (*linodego.NodeBalancerNodeResponse, error)
		listNFn   func(int, int) (*linodego.NodeBalancerNodeListResponse, error)

		service     *v1.Service
		nodes       []*v1.Node
		clusterName string
		nbIP        string
		err         error
	}{
		{
			"update load balancer",
			func(i int, s string, strings map[string]string) (*linodego.NodeBalancerResponse, error) {
				return &linodego.NodeBalancerResponse{
					newFakeOKResponse("nodebalancer.create"),
					linodego.LinodeNodeBalancerId{
						12345,
					},
				}, nil
			},
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				label := linodego.CustomString{}
				err := label.UnmarshalJSON([]byte("afoobar123"))
				if err != nil {
					return nil, err
				}
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{
						{
							// loadbalancer names are a + service.UID
							// see cloudprovider.GetLoadBalancerName
							NodeBalancerId: 12345,
							Label:          label,
							Address4:       "127.0.0.1",
						},
					},
				}, nil
			},
			func(i int, strings map[string]string) (*linodego.NodeBalancerConfigResponse, error) {
				port := strings["Port"]
				if port == "80" || port == "8443" {
					return nil, fmt.Errorf("config already exists")
				}
				return &linodego.NodeBalancerConfigResponse{
					newFakeOKResponse("nodebalancer.config.create"),
					linodego.NodeBalancerConfigId{125},
				}, nil
			},
			func(i int, i2 int) (*linodego.NodeBalancerConfigListResponse, error) {
				return &linodego.NodeBalancerConfigListResponse{
					newFakeOKResponse("nodebalancer.config.list"),
					[]linodego.NodeBalancerConfig{
						{
							Port:     80,
							ConfigId: 123,
						},
						{
							Port:     8443,
							ConfigId: 124,
						},
					},
				}, nil
			},
			func(i int, i2 int) (*linodego.NodeBalancerConfigResponse, error) {
				return &linodego.NodeBalancerConfigResponse{
					newFakeOKResponse("nodebalancer.config.delete"),
					linodego.NodeBalancerConfigId{i2},
				}, nil
			},
			func(i int, strings map[string]string) (*linodego.NodeBalancerConfigResponse, error) {
				return &linodego.NodeBalancerConfigResponse{
					newFakeOKResponse("nodebalancer.config.update"),
					linodego.NodeBalancerConfigId{i},
				}, nil
			},

			func(i int) (*linodego.NodeBalancerNodeResponse, error) {
				return &linodego.NodeBalancerNodeResponse{
					newFakeOKResponse("nodebalancer.node"),
					linodego.NodeBalancerNodeId{i},
				}, nil
			},
			func(i int, strings map[string]string) (*linodego.NodeBalancerNodeResponse, error) {
				if len(strings) == 0 {
					return nil, fmt.Errorf("nothing to update")
				}
				return &linodego.NodeBalancerNodeResponse{
					newFakeOKResponse("nodebalancer.node.update"),
					linodego.NodeBalancerNodeId{i},
				}, nil
			},
			func(i int, s string, s2 string, strings map[string]string) (*linodego.NodeBalancerNodeResponse, error) {
				/*if s == "node-1" || s == "node-4" {
					return nil, fmt.Errorf("node already exists")
				}*/
				return &linodego.NodeBalancerNodeResponse{
					newFakeOKResponse("nodebalancer.node.create"),
					linodego.NodeBalancerNodeId{3},
				}, nil
			},

			func(int, int) (*linodego.NodeBalancerNodeListResponse, error) {
				n1 := linodego.CustomString{}
				if err := n1.UnmarshalJSON([]byte("node-1")); err != nil {
					return nil, err
				}
				n2 := linodego.CustomString{}
				if err := n2.UnmarshalJSON([]byte("node-4")); err != nil {
					return nil, err
				}
				return &linodego.NodeBalancerNodeListResponse{
					newFakeOKResponse("nodebalancer.node.list"),
					[]linodego.NodeBalancerNode{
						{
							NodeId: 1,
							Label:  n1,
						},
						{
							NodeId: 2,
							Label:  n2,
						},
					},
				}, nil
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(8000),
							NodePort: int32(30000),
						},
						{
							Name:     "test2",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30001),
						},
					},
				},
			},
			[]*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-1",
					},
					Status: v1.NodeStatus{
						Addresses: []v1.NodeAddress{
							{
								v1.NodeInternalIP,
								"127.0.0.1",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-2",
					},
					Status: v1.NodeStatus{
						Addresses: []v1.NodeAddress{
							{
								v1.NodeInternalIP,
								"127.0.0.2",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-3",
					},
					Status: v1.NodeStatus{
						Addresses: []v1.NodeAddress{
							{
								v1.NodeInternalIP,
								"127.0.0.3",
							},
						},
					},
				},
			},
			"linodelb",
			"127.0.0.1",
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeNB := &fakeNodeBalancerService{
				listFn: test.listNBFn,
			}
			fakeNBC := &fakeNodeBalancerConfigService{
				createFn: test.createNBCFn,
				listFn:   test.listNBCFn,
				updateFn: test.updateNBCFn,
				deleteFn: test.deleteNBCFn,
			}
			fakeN := &fakeNodeBalancerNodeService{
				createFn: test.createNFn,
				listFn:   test.listNFn,
				updateFn: test.updateNFn,
				deleteFn: test.deleteNFn,
			}

			fakeClient := newFakeLBClient(fakeNB, nil, fakeNBC, fakeN)
			lb := &loadbalancers{fakeClient, "1"}
			lbStatus, err := lb.EnsureLoadBalancer(context.Background(), test.clusterName, test.service, test.nodes)
			if lbStatus.Ingress[0].IP != test.nbIP {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.nbIP)
				t.Logf("actual: %v", lbStatus.Ingress)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_GetLoadBalancer(t *testing.T) {
	testcases := []struct {
		name        string
		listNBFn    func(int) (*linodego.NodeBalancerListResponse, error)
		service     *v1.Service
		clusterName string
		found       bool
		err         error
	}{
		{
			"Load balancer exists",
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				label := linodego.CustomString{}
				err := label.UnmarshalJSON([]byte("afoobar123"))
				if err != nil {
					return nil, err
				}
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{
						{
							// loadbalancer names are a + service.UID
							// see cloudprovider.GetLoadBalancerName
							NodeBalancerId: 12345,
							Label:          label,
							Address4:       "127.0.0.1",
						},
					},
				}, nil
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			"linodelb",
			true,
			nil,
		},
		{
			"Load balancer not exists",
			func(i int) (*linodego.NodeBalancerListResponse, error) {
				label := linodego.CustomString{}
				err := label.UnmarshalJSON([]byte("afoobar321"))
				if err != nil {
					return nil, err
				}
				return &linodego.NodeBalancerListResponse{
					newFakeOKResponse("nodebalancer.list"),
					[]linodego.LinodeNodeBalancer{
						{
							// loadbalancer names are a + service.UID
							// see cloudprovider.GetLoadBalancerName
							NodeBalancerId: 12345,
							Label:          label,
							Address4:       "127.0.0.1",
						},
					},
				}, nil
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "foobar123",
					Annotations: map[string]string{
						annLinodeProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			"linodelb",
			false,
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeNB := &fakeNodeBalancerService{
				listFn: test.listNBFn,
			}

			fakeClient := newFakeLBClient(fakeNB, nil, nil, nil)
			lb := &loadbalancers{fakeClient, "1"}
			_, found, err := lb.GetLoadBalancer(context.Background(), test.clusterName, test.service)
			if found != test.found {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.found)
				t.Logf("actual: %v", found)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}
