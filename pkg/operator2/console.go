package operator2

import (
	"net/url"
	"regexp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	configv1 "github.com/openshift/api/config/v1"
)

func (c *authOperator) handleConsoleConfig() *configv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	consoleConfig, err := c.console.Get(globalConfigName, metav1.GetOptions{})
	if err != nil {
		klog.Infof("error getting console config: %v", err)
		return &configv1.Console{}
	}
	return consoleConfig
}
func consoleToDeploymentData(console *configv1.Console) (string, []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	assetPublicURL := console.Status.ConsoleURL
	if len(assetPublicURL) == 0 {
		return "", nil
	}
	corsAllowedOrigins := []string{"^" + regexp.QuoteMeta(assetPublicURL) + "$"}
	if _, err := url.Parse(assetPublicURL); err != nil {
		klog.Errorf("failed to parse assetPublicURL %s: %v", assetPublicURL, err)
		return "", nil
	}
	for _, corsAllowedOrigin := range corsAllowedOrigins {
		if _, err := regexp.Compile(corsAllowedOrigin); err != nil {
			klog.Errorf("failed to parse corsAllowedOrigin %s: %v", corsAllowedOrigin, err)
			return "", nil
		}
	}
	return assetPublicURL, nil
}
