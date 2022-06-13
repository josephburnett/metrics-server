module sigs.k8s.io/metrics-server

go 1.16

require (
	github.com/go-openapi/spec v0.20.3
	github.com/google/addlicense v0.0.0-20210428195630-6d92264d7170
	github.com/google/go-cmp v0.5.5
	github.com/mailru/easyjson v0.7.7
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.11.0
	github.com/prometheus/common v0.25.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	golang.org/x/perf v0.0.0-20210220033136-40a54f11e909
	k8s.io/api v0.21.5
	k8s.io/apimachinery v0.21.5
	k8s.io/apiserver v0.21.5
	k8s.io/client-go v0.21.5
	k8s.io/component-base v0.21.5
	k8s.io/controller-manager v0.21.5 // indirect
	k8s.io/klog/hack/tools v0.0.0-20210512110738-02ca14bed863
	k8s.io/klog/v2 v2.8.0
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
	k8s.io/kubelet v0.21.5
	k8s.io/metrics v0.21.5
	sigs.k8s.io/mdtoc v1.0.1
)

replace (
        k8s.io/controller-manager/pkg/clientbuilder => k8s.io/controller-manager/pkg/clientbuilder v0.21.5
        k8s.io/api => k8s.io/api v0.21.5
        k8s.io/apiextensions => k8s.io/apiextensions v0.21.5
        k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.5
        k8s.io/apimachinery => k8s.io/apimachinery v0.21.5
        k8s.io/apiserver => k8s.io/apiserver v0.21.5
        k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.5
        k8s.io/client-go => k8s.io/client-go v0.21.5
        k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.5
        k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.5
        k8s.io/code-generator => k8s.io/code-generator v0.21.5
        k8s.io/component-base => k8s.io/component-base v0.21.5
        k8s.io/component-helpers => k8s.io/component-helpers v0.21.5
        k8s.io/controller-manager => k8s.io/controller-manager v0.21.5
        k8s.io/cri-api => k8s.io/cri-api v0.21.5
        k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.5
        k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.5
        k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.5
        k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.5
)
