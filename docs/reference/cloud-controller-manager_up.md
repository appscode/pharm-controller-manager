## cloud-controller-manager up

Bootstrap as a Kubernetes master or node

### Synopsis

Bootstrap as a Kubernetes master or node

```
cloud-controller-manager up [flags]
```

### Options

```
      --address ip                              DEPRECATED: the IP address on which to listen for the --port port (set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces). See --bind-address instead. (default 0.0.0.0)
      --allocate-node-cidrs                     Should CIDRs for Pods be allocated and set on the cloud provider.
      --bind-address ip                         The IP address on which to listen for the --secure-port port. The associated interface(s) must be reachable by the rest of the cluster, and by CLI/web clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces). (default 0.0.0.0)
      --cert-dir string                         The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored. (default "/var/run/kubernetes")
      --cidr-allocator-type string              Type of CIDR allocator to use (default "RangeAllocator")
      --cloud-config string                     The path to the cloud provider configuration file. Empty string for no configuration file.
      --cloud-provider string                   The provider for cloud services. Empty string for no provider.
      --cluster-cidr string                     CIDR Range for Pods in cluster. Requires --allocate-node-cidrs to be true
      --cluster-name string                     The instance prefix for the cluster. (default "kubernetes")
      --concurrent-service-syncs int32          The number of services that are allowed to sync concurrently. Larger number = more responsive service management, but more CPU (and network) load (default 1)
      --configure-cloud-routes                  Should CIDRs allocated by allocate-node-cidrs be configured on the cloud provider. (default true)
      --contention-profiling                    Enable lock contention profiling, if profiling is enabled
      --controller-start-interval duration      Interval between starting controller managers.
      --feature-gates mapStringBool             A set of key=value pairs that describe feature gates for alpha/experimental features. Options are:
                                                APIListChunking=true|false (BETA - default=true)
                                                APIResponseCompression=true|false (ALPHA - default=false)
                                                AdvancedAuditing=true|false (BETA - default=true)
                                                AllAlpha=true|false (ALPHA - default=false)
                                                AppArmor=true|false (BETA - default=true)
                                                AttachVolumeLimit=true|false (ALPHA - default=false)
                                                BalanceAttachedNodeVolumes=true|false (ALPHA - default=false)
                                                BlockVolume=true|false (ALPHA - default=false)
                                                CPUManager=true|false (BETA - default=true)
                                                CRIContainerLogRotation=true|false (BETA - default=true)
                                                CSIBlockVolume=true|false (ALPHA - default=false)
                                                CSIPersistentVolume=true|false (BETA - default=true)
                                                CustomPodDNS=true|false (BETA - default=true)
                                                CustomResourceSubresources=true|false (BETA - default=true)
                                                CustomResourceValidation=true|false (BETA - default=true)
                                                DebugContainers=true|false (ALPHA - default=false)
                                                DevicePlugins=true|false (BETA - default=true)
                                                DynamicKubeletConfig=true|false (BETA - default=true)
                                                DynamicProvisioningScheduling=true|false (ALPHA - default=false)
                                                EnableEquivalenceClassCache=true|false (ALPHA - default=false)
                                                ExpandInUsePersistentVolumes=true|false (ALPHA - default=false)
                                                ExpandPersistentVolumes=true|false (BETA - default=true)
                                                ExperimentalCriticalPodAnnotation=true|false (ALPHA - default=false)
                                                ExperimentalHostUserNamespaceDefaulting=true|false (BETA - default=false)
                                                GCERegionalPersistentDisk=true|false (BETA - default=true)
                                                HugePages=true|false (BETA - default=true)
                                                HyperVContainer=true|false (ALPHA - default=false)
                                                Initializers=true|false (ALPHA - default=false)
                                                KubeletPluginsWatcher=true|false (ALPHA - default=false)
                                                LocalStorageCapacityIsolation=true|false (BETA - default=true)
                                                MountContainers=true|false (ALPHA - default=false)
                                                MountPropagation=true|false (BETA - default=true)
                                                PersistentLocalVolumes=true|false (BETA - default=true)
                                                PodPriority=true|false (BETA - default=true)
                                                PodReadinessGates=true|false (BETA - default=false)
                                                PodShareProcessNamespace=true|false (ALPHA - default=false)
                                                QOSReserved=true|false (ALPHA - default=false)
                                                ReadOnlyAPIDataVolumes=true|false (DEPRECATED - default=true)
                                                ResourceLimitsPriorityFunction=true|false (ALPHA - default=false)
                                                ResourceQuotaScopeSelectors=true|false (ALPHA - default=false)
                                                RotateKubeletClientCertificate=true|false (BETA - default=true)
                                                RotateKubeletServerCertificate=true|false (ALPHA - default=false)
                                                RunAsGroup=true|false (ALPHA - default=false)
                                                ScheduleDaemonSetPods=true|false (ALPHA - default=false)
                                                ServiceNodeExclusion=true|false (ALPHA - default=false)
                                                ServiceProxyAllowExternalIPs=true|false (DEPRECATED - default=false)
                                                StorageObjectInUseProtection=true|false (default=true)
                                                StreamingProxyRedirects=true|false (BETA - default=true)
                                                SupportIPVSProxyMode=true|false (default=true)
                                                SupportPodPidsLimit=true|false (ALPHA - default=false)
                                                Sysctls=true|false (BETA - default=true)
                                                TaintBasedEvictions=true|false (ALPHA - default=false)
                                                TaintNodesByCondition=true|false (ALPHA - default=false)
                                                TokenRequest=true|false (ALPHA - default=false)
                                                TokenRequestProjection=true|false (ALPHA - default=false)
                                                VolumeScheduling=true|false (BETA - default=true)
                                                VolumeSubpath=true|false (default=true)
                                                VolumeSubpathEnvExpansion=true|false (ALPHA - default=false)
  -h, --help                                    help for up
      --http2-max-streams-per-connection int    The limit that the server gives to clients for the maximum number of streams in an HTTP/2 connection. Zero means to use golang's default.
      --kube-api-burst int32                    Burst to use while talking with kubernetes apiserver. (default 30)
      --kube-api-content-type string            Content type of requests sent to apiserver. (default "application/vnd.kubernetes.protobuf")
      --kube-api-qps float32                    QPS to use while talking with kubernetes apiserver. (default 20)
      --kubeconfig string                       Path to kubeconfig file with authorization and master location information.
      --leader-elect                            Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability. (default true)
      --leader-elect-lease-duration duration    The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled. (default 15s)
      --leader-elect-renew-deadline duration    The interval between attempts by the acting master to renew a leadership slot before it stops leading. This must be less than or equal to the lease duration. This is only applicable if leader election is enabled. (default 10s)
      --leader-elect-resource-lock endpoints    The type of resource object that is used for locking during leader election. Supported options are endpoints (default) and `configmaps`. (default "endpoints")
      --leader-elect-retry-period duration      The duration the clients should wait between attempting acquisition and renewal of a leadership. This is only applicable if leader election is enabled. (default 2s)
      --master string                           The address of the Kubernetes API server (overrides any value in kubeconfig).
      --min-resync-period duration              The resync period in reflectors will be random between MinResyncPeriod and 2*MinResyncPeriod. (default 12h0m0s)
      --node-monitor-period duration            The period for syncing NodeStatus in NodeController. (default 5s)
      --node-status-update-frequency duration   Specifies how often the controller updates nodes' status. (default 5m0s)
      --port int                                DEPRECATED: the port on which to serve HTTP insecurely without authentication and authorization. If 0, don't serve HTTPS at all. See --secure-port instead. (default 10253)
      --profiling                               Enable profiling via web interface host:port/debug/pprof/
      --route-reconciliation-period duration    The period for reconciling routes created for Nodes by cloud provider. (default 10s)
      --secure-port int                         The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.
      --tls-cert-file string                    File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert). If HTTPS serving is enabled, and --tls-cert-file and --tls-private-key-file are not provided, a self-signed certificate and key are generated for the public address and saved to the directory specified by --cert-dir.
      --tls-cipher-suites strings               Comma-separated list of cipher suites for the server. If omitted, the default Go cipher suites will be use.  Possible values: TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_RC4_128_SHA,TLS_RSA_WITH_3DES_EDE_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA256,TLS_RSA_WITH_AES_128_GCM_SHA256,TLS_RSA_WITH_AES_256_CBC_SHA,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_RC4_128_SHA
      --tls-min-version string                  Minimum TLS version supported. Possible values: VersionTLS10, VersionTLS11, VersionTLS12
      --tls-private-key-file string             File containing the default x509 private key matching --tls-cert-file.
      --tls-sni-cert-key namedCertKey           A pair of x509 certificate and private key file paths, optionally suffixed with a list of domain patterns which are fully qualified domain names, possibly with prefixed wildcard segments. If no domain patterns are provided, the names of the certificate are extracted. Non-wildcard matches trump over wildcard matches, explicit domain patterns trump over extracted names. For multiple key/certificate pairs, use the --tls-sni-cert-key multiple times. Examples: "example.crt,example.key" or "foo.crt,foo.key:*.foo.com,foo.com". (default [])
      --use-service-account-credentials         If true, use individual service account credentials for each controller.
```

### Options inherited from parent commands

```
      --alsologtostderr                         log to standard error as well as files
      --analytics                               Send analytical events to Google Analytics (default true)
      --cloud-provider-gce-lb-src-cidrs cidrs   CIDRs opened in GCE firewall for LB traffic proxy & health checks (default 130.211.0.0/22,209.85.152.0/22,209.85.204.0/22,35.191.0.0/16)
      --log_backtrace_at traceLocation          when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                          If non-empty, write log files in this directory
      --logtostderr                             log to standard error instead of files
      --stderrthreshold severity                logs at or above this threshold go to stderr (default 2)
  -v, --v Level                                 log level for V logs
      --version version[=true]                  Print version information and quit
      --vmodule moduleSpec                      comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [cloud-controller-manager](cloud-controller-manager.md)	 - Pharm Controller Manager by Appscode - Start farms

