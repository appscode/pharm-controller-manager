package cmds

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/appscode/go/log"
	"github.com/golang/glog"
	_ "github.com/pharmer/cloud-controller-manager/cloud/providers"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	"k8s.io/kubernetes/pkg/cloudprovider"
	_ "k8s.io/kubernetes/pkg/version/prometheus" // for version metric registration
)

func init() {
	healthz.DefaultHealthz()
}

func NewCmdUp() *cobra.Command {
	s := options.NewCloudControllerManagerServer()
	cmd := &cobra.Command{
		Use:               "up",
		Short:             "Bootstrap as a Kubernetes master or node",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := wait.Poll(3*time.Second, 5*time.Minute, func() (bool, error) {
				addr, err := net.LookupIP("google.com")
				return len(addr) > 0, err
			})
			if err != nil {
				log.Fatalln("Failed to resolve DNS. Reason: %v", err)
			}

			fmt.Println(s.CloudProvider, "*****************")

			cloud, err := cloudprovider.InitCloudProvider(s.CloudProvider, s.CloudConfigFile)
			fmt.Println(s.CloudConfigFile, "----")
			if err != nil {
				glog.Fatalf("Cloud provider could not be initialized: %v", err)
			}

			if err := app.Run(s, cloud); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	s.AddFlags(cmd.Flags())
	return cmd
}
