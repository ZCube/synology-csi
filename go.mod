module github.com/SynologyOpenSource/synology-csi

go 1.20

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/container-storage-interface/spec v1.7.0
	github.com/kubernetes-csi/csi-lib-utils v0.9.1
	github.com/kubernetes-csi/csi-test/v4 v4.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.6.0
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.26.1
	k8s.io/apimachinery v0.26.1 // indirect
	k8s.io/client-go v0.26.1 // indirect
	k8s.io/cloud-provider v0.26.1 // indirect
	k8s.io/component-base v0.26.1 // indirect
	k8s.io/klog/v2 v2.90.0 // indirect
	k8s.io/kubernetes v1.26.0
	k8s.io/mount-utils v0.26.1
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448
)

exclude google.golang.org/grpc v1.37.0

require (
	github.com/satori/uuid v1.2.0
	k8s.io/apiserver v0.26.1 // indirect
	k8s.io/component-helpers v0.26.1 // indirect
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.23.0 // indirect
	github.com/opencontainers/selinux v1.10.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/net v0.3.1-0.20221206200815-1e63c2f08a10 // indirect
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
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
)
