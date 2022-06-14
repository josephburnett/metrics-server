package bridge

import (
	"context"
	"fmt"
	"sync"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

// The in-memory pod resource metrics client should be injected into
// the HPA controller. The singleton bridge is a hacky shortcut.
var bridgePodResourceMetricsesGetter *podResourceMetricsesGetter = &podResourceMetricsesGetter{}

type podResourceMetricsesGetter struct {
	mux  sync.Mutex
	list func(ctx context.Context, opts *metainternalversion.ListOptions) (runtime.Object, error)
}

func (p *podResourceMetricsesGetter) PodMetricses(namespace string) resourceclient.PodMetricsInterface {
	p.mux.Lock()
	defer p.mux.Unlock()
	return &podResourceMetrics{
		namespace: namespace,
		list:      p.list,
	}
}

func GetPodResourceMetricsClient() resourceclient.PodMetricsesGetter {
	return bridgePodResourceMetricsesGetter
}

func SetPodResourceMetricsListFn(list func(ctx context.Context, opts *metainternalversion.ListOptions) (runtime.Object, error)) {
	b := bridgePodResourceMetricsesGetter
	b.mux.Lock()
	defer b.mux.Unlock()
	b.list = list
}

type podResourceMetrics struct {
	namespace string
	list      func(ctx context.Context, opts *metainternalversion.ListOptions) (runtime.Object, error)
}

func (p *podResourceMetrics) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.PodMetrics, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (p *podResourceMetrics) List(parent context.Context, opts v1.ListOptions) (*v1beta1.PodMetricsList, error) {
	ctx := request.WithNamespace(parent, p.namespace)
	if p.list == nil {
		return nil, fmt.Errorf("implementation not provided")
	}
	obj, err := p.list(ctx, optionsAdapter(opts))
	return obj.(*v1beta1.PodMetricsList), err
}

func optionsAdapter(opts v1.ListOptions) *metainternalversion.ListOptions {
	o := &metainternalversion.ListOptions{}
	o.APIVersion = opts.APIVersion
	o.AllowWatchBookmarks = opts.AllowWatchBookmarks
	o.Continue = opts.Continue
	o.FieldSelector = fields.ParseSelectorOrDie(opts.FieldSelector)
	o.Kind = opts.Kind
	o.LabelSelector, _ = labels.Parse(opts.LabelSelector)
	o.Limit = opts.Limit
	o.ResourceVersion = opts.ResourceVersion
	o.ResourceVersionMatch = opts.ResourceVersionMatch
	o.TimeoutSeconds = opts.TimeoutSeconds
	o.TypeMeta = opts.TypeMeta
	o.Watch = opts.Watch
	return o
}

func (p *podResourceMetrics) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("unimplemented")
}
