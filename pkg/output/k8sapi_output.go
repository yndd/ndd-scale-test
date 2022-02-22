package output

import (
	"bytes"
	"context"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sApiOutput struct {
	client       dynamic.Interface
	mapper       *restmapper.DeferredDiscoveryRESTMapper
	committedCRs []*CommittedCR
	kubeconfig   string
}

type CommittedCR struct {
	data    *unstructured.Unstructured
	index   int
	mapping *meta.RESTMapping
}

func NewK8sOutput(kubeconfig string) *K8sApiOutput {
	result := &K8sApiOutput{
		client:       nil,
		mapper:       nil,
		committedCRs: []*CommittedCR{},
	}
	if kubeconfig != "" {
		result.kubeconfig = kubeconfig
	}
	return result
}

func (k *K8sApiOutput) Commit(data *bytes.Buffer, opi *OutputPluginInfo) error {

	err := k.init()
	if err != nil {
		return err
	}

	apiinfo, err := k.k8sCreate(data)
	if err != nil {
		return err
	}

	k.committedCRs = append(k.committedCRs, &CommittedCR{
		data:    apiinfo.obj,
		index:   opi.Index,
		mapping: apiinfo.restmapping,
	})

	return nil
}

func (k *K8sApiOutput) Delete(data *bytes.Buffer, opi *OutputPluginInfo) error {
	err := k.init()
	if err != nil {
		return err
	}

	err = k.k8sDelete(data)
	if err != nil {
		return err
	}
	return nil
}

func deduceApiInformation(b *bytes.Buffer, mapper *restmapper.DeferredDiscoveryRESTMapper) (*ApiInfos, error) {
	obj := &unstructured.Unstructured{}

	// decode YAML into unstructured.Unstructured
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode(b.Bytes(), nil, obj)
	if err != nil {
		return nil, err
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	namespaceValue := obj.GetNamespace()

	return &ApiInfos{namespace: namespaceValue, restmapping: mapping, obj: obj}, nil
}

type ApiInfos struct {
	namespace   string
	restmapping *meta.RESTMapping
	obj         *unstructured.Unstructured
}

func (k *K8sApiOutput) k8sCreate(b *bytes.Buffer) (*ApiInfos, error) {

	apiinfo, err := deduceApiInformation(b, k.mapper)
	if err != nil {
		return nil, err
	}
	log.Infof("Creating %s of kind %s in %s\n", apiinfo.obj.GetName(), apiinfo.obj.GetKind(), apiinfo.namespace)
	_, err = k.client.Resource(apiinfo.restmapping.Resource).Namespace(apiinfo.namespace).Create(context.TODO(), apiinfo.obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return apiinfo, nil
}

func (k *K8sApiOutput) k8sDelete(b *bytes.Buffer) error {
	apiinfo, err := deduceApiInformation(b, k.mapper)
	if err != nil {
		return err
	}

	log.Infof("Creating %s of kind %s in %s\n", apiinfo.obj.GetName(), apiinfo.obj.GetKind(), apiinfo.namespace)

	err = k.client.Resource(apiinfo.restmapping.Resource).Namespace(apiinfo.namespace).Delete(context.TODO(), apiinfo.obj.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// init initializes the kubernetes client
func (k *K8sApiOutput) init() error {

	// initialize if not yet done!
	if k.client != nil {
		return nil
	}
	log.Info("Creating new Kubernetes client")
	kubeconfig := k.kubeconfig
	if kubeconfig == "" {
		home := homedir.HomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	k.client = client

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	k.mapper = mapper
	return nil
}
