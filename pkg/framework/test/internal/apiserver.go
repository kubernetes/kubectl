package internal

import (
	"fmt"
	"net/url"
)

func MakeAPIServerArgs(ps DefaultedProcessInput, etcdURL *url.URL) ([]string, error) {
	if etcdURL == nil {
		return []string{}, fmt.Errorf("must configure Etcd URL")
	}

	args := []string{
		"--authorization-mode=Node,RBAC",
		"--runtime-config=admissionregistration.k8s.io/v1alpha1",
		"--v=3", "--vmodule=",
		"--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,SecurityContextDeny,DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota",
		"--admission-control-config-file=",
		"--bind-address=0.0.0.0",
		"--storage-backend=etcd3",
		fmt.Sprintf("--etcd-servers=%s", etcdURL.String()),
		fmt.Sprintf("--cert-dir=%s", ps.Dir),
		fmt.Sprintf("--insecure-port=%s", ps.URL.Port()),
		fmt.Sprintf("--insecure-bind-address=%s", ps.URL.Hostname()),
	}

	return args, nil
}
