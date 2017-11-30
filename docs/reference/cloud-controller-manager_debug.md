## cloud-controller-manager debug

Debug cloud-controller-manager

### Synopsis


Debug cloud-controller-manager

```
cloud-controller-manager debug [flags]
```

### Options

```
      --cloud-config string   The path to the cloud provider configuration file.  Empty string for no configuration file.
  -h, --help                  help for debug
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

