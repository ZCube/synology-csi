module github.com/SynologyOpenSource/synology-csi

go 1.16

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/container-storage-interface/spec v1.7.0
	github.com/kubernetes-csi/csi-lib-utils v0.9.1
	github.com/kubernetes-csi/csi-test/v4 v4.3.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.3
	golang.org/x/sys v0.5.0
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/mount-utils v0.26.4
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d
	k8s.io/kubernetes v1.26.0
)

exclude google.golang.org/grpc v1.37.0

require (
	k8s.io/apiserver v0.26.1 // indirect
	k8s.io/component-helpers v0.26.1 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.26.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.0
	k8s.io/apiserver => k8s.io/apiserver v0.26.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.26.0
	k8s.io/client-go => k8s.io/client-go v0.26.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.26.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.0
	k8s.io/code-generator => k8s.io/code-generator v0.26.0
	k8s.io/component-base => k8s.io/component-base v0.26.0
	k8s.io/component-helpers => k8s.io/component-helpers v0.26.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.26.0
	k8s.io/cri-api => k8s.io/cri-api v0.26.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.26.0
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.26.0
	k8s.io/kms => k8s.io/kms v0.26.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.26.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.26.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.26.0
	k8s.io/kubectl => k8s.io/kubectl v0.26.0
	k8s.io/kubelet => k8s.io/kubelet v0.26.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.26.0
	k8s.io/metrics => k8s.io/metrics v0.26.0
	k8s.io/mount-utils => k8s.io/mount-utils v0.0.0-20230103133730-1df1a57439e2
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.26.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.26.0
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.26.0
	k8s.io/sample-controller => k8s.io/sample-controller v0.26.0
	k8s.io/kubernetes => k8s.io/kubernetes v0.26.0
)
