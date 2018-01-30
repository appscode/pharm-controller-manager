package linode

import (
	"reflect"
	"testing"

	"github.com/taoh/linodego"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type fakeLinodeService struct {
	bootFunc     func(int, int) (*linodego.JobResponse, error)
	cloneFunc    func(int, int, int, int) (*linodego.LinodeResponse, error)
	createFunc   func(int, int, int) (*linodego.LinodeResponse, error)
	deleteFunc   func(int, bool) (*linodego.LinodeResponse, error)
	listFunc     func(int) (*linodego.LinodesListResponse, error)
	rebootFunc   func(int, int) (*linodego.JobResponse, error)
	resizeFunc   func(int, int) (*linodego.LinodeResponse, error)
	shutdownFunc func(int) (*linodego.JobResponse, error)
	updateFunc   func(int, map[string]interface{}) (*linodego.LinodeResponse, error)
}

func (f *fakeLinodeService) Boot(linodeId int, configId int) (*linodego.JobResponse, error) {
	return f.bootFunc(linodeId, configId)
}

func (f *fakeLinodeService) Clone(linodeId int, dataCenterId int, planId int, paymentTerm int) (*linodego.LinodeResponse, error) {
	return f.cloneFunc(linodeId, dataCenterId, planId, paymentTerm)
}

func (f *fakeLinodeService) Create(dataCenterId int, planId int, paymentTerm int) (*linodego.LinodeResponse, error) {
	return f.createFunc(dataCenterId, planId, paymentTerm)
}

func (f *fakeLinodeService) Delete(linodeId int, skipChecks bool) (*linodego.LinodeResponse, error) {
	return f.deleteFunc(linodeId, skipChecks)
}

func (f *fakeLinodeService) List(linodeId int) (*linodego.LinodesListResponse, error) {
	return f.listFunc(linodeId)
}

func (f *fakeLinodeService) Reboot(linodeId int, configId int) (*linodego.JobResponse, error) {
	return f.rebootFunc(linodeId, configId)
}

func (f *fakeLinodeService) Resize(linodeId int, planId int) (*linodego.LinodeResponse, error) {
	return f.resizeFunc(linodeId, planId)
}

func (f *fakeLinodeService) Shutdown(linodeId int) (*linodego.JobResponse, error) {
	return f.shutdownFunc(linodeId)
}

func (f *fakeLinodeService) Update(linodeId int, args map[string]interface{}) (*linodego.LinodeResponse, error) {
	return f.updateFunc(linodeId, args)
}

type fakeLinodeIPService struct {
	addPrivateFn func(int) (*linodego.LinodeIPAddressResponse, error)
	addPublicFn  func(int) (*linodego.LinodeIPAddressResponse, error)
	listFn       func(int, int) (*linodego.LinodeIPListResponse, error)
	setRDNSFn    func(int, string) (*linodego.LinodeRDNSIPAddressResponse, error)
	swapFn       func(int, int, int) (*linodego.LinodeLinodeIPAddressResponse, error)
}

func (f *fakeLinodeIPService) List(linodeId int, ipAddressId int) (*linodego.LinodeIPListResponse, error) {
	return f.listFn(linodeId, ipAddressId)
}

func (f *fakeLinodeIPService) AddPrivate(linodeId int) (*linodego.LinodeIPAddressResponse, error) {
	return f.addPrivateFn(linodeId)
}

func (f *fakeLinodeIPService) AddPublic(linodeId int) (*linodego.LinodeIPAddressResponse, error) {
	return f.addPublicFn(linodeId)
}

func (f *fakeLinodeIPService) SetRDNS(ipAddressId int, hostname string) (*linodego.LinodeRDNSIPAddressResponse, error) {
	return f.setRDNSFn(ipAddressId, hostname)
}

func (f *fakeLinodeIPService) Swap(ipAddressId int, withIPAddressId int, toLinodeId int) (*linodego.LinodeLinodeIPAddressResponse, error) {
	return f.swapFn(ipAddressId, withIPAddressId, toLinodeId)
}
func newFakeClient(fakeL *fakeLinodeService, fakeIP *fakeLinodeIPService) *linodego.Client {
	client := linodego.NewClient("", nil)
	client.Linode = fakeL
	client.Ip = fakeIP
	return client
}

func newFakeOKResponse(action string) linodego.Response {
	return linodego.Response{
		Errors: nil,
		Action: action,
	}
}

func newFakeNotOKResponse(code int, message, action string) linodego.Response {
	return linodego.Response{
		Errors: []linodego.Error{
			{
				code,
				message,
			},
		},
		Action: action,
	}
}

var _ cloudprovider.Instances = new(instances)

func newFakeLinode() linodego.Linode {
	label := linodego.CustomString{}
	err := label.UnmarshalJSON([]byte("test-linode"))
	if err != nil {
		return linodego.Linode{}
	}
	return linodego.Linode{
		Label:    label,
		LinodeId: 1234,
		PlanId:   2,
	}
}

func Test_NodeAddresses(t *testing.T) {
	fake := &fakeLinodeService{}
	fake.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}
	fakeIP := &fakeLinodeIPService{}
	fakeIP.listFn = func(i int, i2 int) (*linodego.LinodeIPListResponse, error) {
		linode := newFakeLinode()
		return &linodego.LinodeIPListResponse{
			newFakeOKResponse("linode.ip.list"),
			[]linodego.FullIPAddress{
				{
					LinodeId:    linode.LinodeId,
					IsPublic:    0,
					IPAddress:   "10.0.0.1",
					IPAddressId: 1,
				},
				{
					LinodeId:    linode.LinodeId,
					IsPublic:    1,
					IPAddress:   "93.58.178.92",
					IPAddressId: 2,
				},
			},
		}, nil
	}
	expectedAddresses := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "test-linode",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.0.0.1",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "93.58.178.92",
		},
	}
	fakeClient := newFakeClient(fake, fakeIP)
	instances := newInstances(fakeClient)

	addresses, err := instances.NodeAddresses("test-linode")
	if !reflect.DeepEqual(addresses, expectedAddresses) {
		t.Errorf("unexpected node addresses. got: %v want: %v", addresses, expectedAddresses)
	}

	if err != nil {
		t.Errorf("unexpected err, expected nil. got: %v", err)
	}

}

func Test_NodeAddressesByProviderID(t *testing.T) {
	fake := &fakeLinodeService{}
	fake.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}
	fakeIP := &fakeLinodeIPService{}
	fakeIP.listFn = func(i int, i2 int) (*linodego.LinodeIPListResponse, error) {
		linode := newFakeLinode()
		return &linodego.LinodeIPListResponse{
			newFakeOKResponse("linode.ip.list"),
			[]linodego.FullIPAddress{
				{
					LinodeId:    linode.LinodeId,
					IsPublic:    0,
					IPAddress:   "10.0.0.1",
					IPAddressId: 1,
				},
				{
					LinodeId:    linode.LinodeId,
					IsPublic:    1,
					IPAddress:   "93.58.178.92",
					IPAddressId: 2,
				},
			},
		}, nil
	}

	expectedAddresses := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "test-linode",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.0.0.1",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "93.58.178.92",
		},
	}
	fakeClient := newFakeClient(fake, fakeIP)
	instances := newInstances(fakeClient)

	addresses, err := instances.NodeAddressesByProviderID("linode://1234")
	if !reflect.DeepEqual(addresses, expectedAddresses) {
		t.Errorf("unexpected node addresses. got: %v want: %v", addresses, expectedAddresses)
	}

	if err != nil {
		t.Errorf("unexpected err, expected nil. got: %v", err)
	}
}

func Test_ExternalID(t *testing.T) {
	fake := &fakeLinodeService{}
	fake.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}
	fakeClient := newFakeClient(fake, nil)
	instances := newInstances(fakeClient)
	id, err := instances.ExternalID("test-linode")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if id != "1234" {
		t.Errorf("expected id 1234, got %s", id)
	}

}

func Test_InstanceType(t *testing.T) {
	fake := &fakeLinodeService{}
	fake.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}
	fakeClient := newFakeClient(fake, nil)
	instances := newInstances(fakeClient)

	insType, err := instances.InstanceType("test-linode")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if insType != "2" {
		t.Errorf("expected id 2, got %s", insType)
	}
}

func Test_InstanceTypeByProviderID(t *testing.T) {
	fake := &fakeLinodeService{}
	fake.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}
	fakeClient := newFakeClient(fake, nil)
	instances := newInstances(fakeClient)

	insType, err := instances.InstanceTypeByProviderID("linode://1234")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if insType != "2" {
		t.Errorf("expected id 2, got %s", insType)
	}
}

func Test_InstanceExistsByProviderID(t *testing.T) {
	fake := &fakeLinodeService{}
	fake.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		if i == linode.LinodeId {
			return &linodego.LinodesListResponse{
				newFakeOKResponse("linode.list"),
				linodes,
			}, nil
		}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			[]linodego.Linode{},
		}, nil
	}
	fakeClient := newFakeClient(fake, nil)
	instances := newInstances(fakeClient)

	found, err := instances.InstanceExistsByProviderID("linode://1234")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if !found {
		t.Errorf("expected found true, got %v", found)
	}

	found, err = instances.InstanceExistsByProviderID("linode://12345")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if found {
		t.Errorf("expected found false, got %v", found)
	}

}
