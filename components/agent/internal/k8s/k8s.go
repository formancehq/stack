package k8s

import (
	"context"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	clientv1beta3 "github.com/formancehq/operator/pkg/client/v1beta3"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func newClient(config *rest.Config) (*clientv1beta3.Client, error) {
	sharedlogging.FromContext(context.Background()).Infof("Connect to cluster")
	defer func() {
		sharedlogging.FromContext(context.Background()).Infof("Connect to cluster OK")
	}()
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &v1beta3.GroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	if crdConfig.UserAgent == "" {
		crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	client, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}

	return clientv1beta3.NewClient(client), nil
}
