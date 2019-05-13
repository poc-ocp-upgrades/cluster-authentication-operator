package operator2

import (
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func (c *authOperator) handleInfrastructureConfig() *configv1.Infrastructure {
	_logClusterCodePath()
	defer _logClusterCodePath()
	infrastructureConfig, err := c.infrastructure.Get(globalConfigName, metav1.GetOptions{})
	if err != nil {
		klog.Infof("error getting infrastructure config: %v", err)
		return &configv1.Infrastructure{Status: configv1.InfrastructureStatus{APIServerURL: "<api_server_url>"}}
	}
	return infrastructureConfig
}
