package bridge

import (
	"context"
	"fmt"
	"sync"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/metrics/pkg/apis/metrics"
)

// The in-memory pod resource metrics client should be injected into
// the HPA controller. The singleton bridge is a hacky shortcut.
var bridgePodResourceMetrics = &podResourceMetrics{}

type podResourceMetrics struct {
	mux  sync.Mutex
	list func(ctx context.Context, opts *metainternalversion.ListOptions) (runtime.Object, error)
}

func List(namespace string, selector labels.Selector) (*metrics.PodMetricsList, error) {
	p := bridgePodResourceMetrics
	p.mux.Lock()
	defer p.mux.Unlock()
	ctx := request.WithNamespace(context.TODO(), namespace)
	if p.list == nil {
		return nil, fmt.Errorf("implementation not provided")
	}
	obj, err := p.list(ctx, &metainternalversion.ListOptions{LabelSelector: selector})
	return obj.(*metrics.PodMetricsList), err
}

func SetPodResourceMetricsListFn(list func(ctx context.Context, opts *metainternalversion.ListOptions) (runtime.Object, error)) {
	b := bridgePodResourceMetrics
	b.mux.Lock()
	defer b.mux.Unlock()
	b.list = list
}
