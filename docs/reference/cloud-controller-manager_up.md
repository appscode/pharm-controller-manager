## cloud-controller-manager up

Bootstrap as a Kubernetes master or node

### Synopsis


Bootstrap as a Kubernetes master or node

```
cloud-controller-manager up [flags]
```

### Options

```
      --address ip                              The IP address to serve on (set to 0.0.0.0 for all interfaces). (default 0.0.0.0)
      --allocate-node-cidrs                     Should CIDRs for Pods be allocated and set on the cloud provider.
      --cloud-config string                     The path to the cloud provider configuration file. Empty string for no configuration file.
      --cloud-provider string                   The provider of cloud services. Cannot be empty.
      --cluster-cidr string                     CIDR Range for Pods in cluster.
      --cluster-name string                     The instance prefix for the cluster. (default "kubernetes")
      --configure-cloud-routes                  Should CIDRs allocated by allocate-node-cidrs be configured on the cloud provider. (default true)
      --contention-profiling                    Enable lock contention profiling, if profiling is enabled.
      --controller-start-interval duration      Interval between starting controller managers.
      --feature-gates mapStringBool             A set of key=value pairs that describe feature gates for alpha/experimental features. Options are:
APIListChunking=true|false (ALPHA - default=false)
APIResponseCompression=true|false (ALPHA - default=false)
Accelerators=true|false (ALPHA - default=false)
AdvancedAuditing=true|false (BETA - default=true)
AllAlpha=true|false (ALPHA - default=false)
AllowExtTrafficLocalEndpoints=true|false (default=true)
AppArmor=true|false (BETA - default=true)
CPUManager=true|false (ALPHA - default=false)
CustomResourceValidation=true|false (ALPHA - default=false)
DebugContainers=true|false (ALPHA - default=false)
DevicePlugins=true|false (ALPHA - default=false)
DynamicKubeletConfig=true|false (ALPHA - default=false)
EnableEquivalenceClassCache=true|false (ALPHA - default=false)
ExpandPersistentVolumes=true|false (ALPHA - default=false)
ExperimentalCriticalPodAnnotation=true|false (ALPHA - default=false)
ExperimentalHostUserNamespaceDefaulting=true|false (BETA - default=false)
HugePages=true|false (ALPHA - default=false)
Initializers=true|false (ALPHA - default=false)
KubeletConfigFile=true|false (ALPHA - default=false)
LocalStorageCapacityIsolation=true|false (ALPHA - default=false)
MountPropagation=true|false (ALPHA - default=false)
PersistentLocalVolumes=true|false (ALPHA - default=false)
PodPriority=true|false (ALPHA - default=false)
RotateKubeletClientCertificate=true|false (BETA - default=true)
RotateKubeletServerCertificate=true|false (ALPHA - default=false)
StreamingProxyRedirects=true|false (BETA - default=true)
SupportIPVSProxyMode=true|false (ALPHA - default=false)
TaintBasedEvictions=true|false (ALPHA - default=false)
TaintNodesByCondition=true|false (ALPHA - default=false)
  -h, --help                                    help for up
      --kube-api-burst int32                    Burst to use while talking with kubernetes apiserver. (default 30)
      --kube-api-content-type string            Content type of requests sent to apiserver. (default "application/vnd.kubernetes.protobuf")
      --kube-api-qps float32                    QPS to use while talking with kubernetes apiserver. (default 20)
      --kubeconfig string                       Path to kubeconfig file with authorization and master location information.
      --leader-elect                            Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability. (default true)
      --leader-elect-lease-duration duration    The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled. (default 15s)
      --leader-elect-renew-deadline duration    The interval between attempts by the acting master to renew a leadership slot before it stops leading. This must be less than or equal to the lease duration. This is only applicable if leader election is enabled. (default 10s)
      --leader-elect-resource-lock endpoints    The type of resource object that is used for locking during leader election. Supported options are endpoints (default) and `configmap`. (default "endpoints")
      --leader-elect-retry-period duration      The duration the clients should wait between attempting acquisition and renewal of a leadership. This is only applicable if leader election is enabled. (default 2s)
      --master string                           The address of the Kubernetes API server (overrides any value in kubeconfig).
      --min-resync-period duration              The resync period in reflectors will be random between MinResyncPeriod and 2*MinResyncPeriod. (default 12h0m0s)
      --node-monitor-period duration            The period for syncing NodeStatus in NodeController. (default 5s)
      --node-status-update-frequency duration   Specifies how often the controller updates nodes' status. (default 5m0s)
      --port int32                              The port that the cloud-controller-manager's http service runs on. (default 10253)
      --profiling                               Enable profiling via web interface host:port/debug/pprof/. (default true)
      --route-reconciliation-period duration    The period for reconciling routes created for Nodes by cloud provider.
      --use-service-account-credentials         If true, use individual service account credentials for each controller.
```

### Options inherited from parent commands

```
      --alsologtostderr                         log to standard error as well as files
      --analytics                               Send analytical events to Google Analytics (default true)
      --cloud-provider-gce-lb-src-cidrs cidrs   CIDRS opened in GCE firewall for LB traffic proxy & health checks (default 35.191.0.0/16,209.85.152.0/22,209.85.204.0/22,130.211.0.0/22)
      --log_backtrace_at traceLocation          when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                          If non-empty, write log files in this directory
      --logtostderr                             log to standard error instead of files
      --stderrthreshold severity                logs at or above this threshold go to stderr (default 2)
  -v, --v Level                                 log level for V logs
      --vmodule moduleSpec                      comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO
* [cloud-controller-manager](cloud-controller-manager.md)	 - Pharm Controller Manager by Appscode - Start farms

