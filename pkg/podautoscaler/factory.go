package podautoscaler

import (
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	kube_client "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	"k8s.io/controller-manager/pkg/clientbuilder"
	resourceclient "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/custom_metrics"
	"k8s.io/metrics/pkg/client/external_metrics"
	horizontalmetrics "sigs.k8s.io/metrics-server/pkg/podautoscaler/metrics"
)

type ControllerFactory struct {
	StopCh     <-chan struct{}
	KubeConfig *restclient.Config
}

func (cf *ControllerFactory) Make() (*HorizontalController, error) {
	kubeClient := kube_client.NewForConfigOrDie(cf.KubeConfig)
	clientBuilder := clientbuilder.SimpleControllerClientBuilder{
		ClientConfig: cf.KubeConfig,
	}
	metricsClientBuilder := clientbuilder.SimpleControllerClientBuilder{
		ClientConfig: cf.KubeConfig,
	}
	factory := informers.NewSharedInformerFactory(kubeClient, 15*time.Minute)

	// Defaults based on:
	//https://github.com/kubernetes/kubernetes/blob/cbdc9b671f33b0f0679e790cc462b25d1476a3af/pkg/controller/apis/config/v1alpha1/defaults.go#L154-L180
	horizontalPodAutoscalerSyncPeriod := 15 * time.Second
	horizontalPodAutoscalerTolerance := 0.1
	horizontalPodAutoscalerCPUInitializationPeriod := 5 * time.Minute
	horizontalPodAutoscalerInitialReadinessDelay := 30 * time.Second
	horizontalPodAutoscalerDownscaleStabilizationWindow := 5 * time.Minute

	// Based on controller-manager construction of rest mapper:
	// https://github.com/kubernetes/kubernetes/blob/6a277e0c4dac8cce4a69cad91fcc7c65de32688c/cmd/kube-controller-manager/app/controllermanager.go#L452-L458
	discoveryClient := clientBuilder.ClientOrDie("controller-discovery")
	cachedClient := cacheddiscovery.NewMemCacheClient(discoveryClient.Discovery())
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)
	go wait.Until(func() {
		restMapper.Reset()
	}, 30*time.Second, cf.StopCh)

	// Based on controller-manager construction of metrics client:
	// https://github.com/kubernetes/kubernetes/blob/99f319567a5148f501e49da35c83478303eab38b/cmd/kube-controller-manager/app/autoscaling.go#L51-L66
	metricsClientConfig := metricsClientBuilder.ConfigOrDie("metrics-horizontal-pod-autoscaler")
	apiVersionsGetter := custom_metrics.NewAvailableAPIsGetter(kubeClient.Discovery())
	// invalidate the discovery information roughly once per resync interval our API
	// information is *at most* two resync intervals old.
	go custom_metrics.PeriodicallyInvalidate(
		apiVersionsGetter,
		horizontalPodAutoscalerSyncPeriod,
		cf.StopCh)
	metricsClient := horizontalmetrics.NewRESTMetricsClient(
		resourceclient.NewForConfigOrDie(metricsClientConfig),
		custom_metrics.NewForConfig(metricsClientConfig, restMapper, apiVersionsGetter),
		external_metrics.NewForConfigOrDie(metricsClientConfig),
	)

	clientConfig := clientBuilder.ConfigOrDie("kuba-horizontal-pod-autoscaler")
	client := clientBuilder.ClientOrDie("kuba-horizontal-pod-autoscaler")
	// Based on controller-manager construction of scale client:
	// https://github.com/kubernetes/kubernetes/blob/99f319567a5148f501e49da35c83478303eab38b/cmd/kube-controller-manager/app/autoscaling.go#L86-L92
	// we don't use cached discovery because DiscoveryScaleKindResolver does its own caching,
	// so we want to re-fetch every time when we actually ask for it
	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(client.Discovery())
	scaleClient, err := scale.NewForConfig(clientConfig, restMapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	if err != nil {
		return nil, err
	}

	hpas := factory.Autoscaling().V2().HorizontalPodAutoscalers()
	go hpas.Informer().Run(cf.StopCh)
	pods := factory.Core().V1().Pods()
	go pods.Informer().Run(cf.StopCh)
	services := factory.Core().V1().Services()
	go services.Informer().Run(cf.StopCh)

	return NewHorizontalController(
		kubeClient.CoreV1(),
		scaleClient,
		kubeClient.AutoscalingV2(),
		restMapper,
		metricsClient,
		hpas,
		pods,
		horizontalPodAutoscalerSyncPeriod,
		horizontalPodAutoscalerDownscaleStabilizationWindow,
		horizontalPodAutoscalerTolerance,
		horizontalPodAutoscalerCPUInitializationPeriod,
		horizontalPodAutoscalerInitialReadinessDelay,
	), nil
}
