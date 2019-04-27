package operator2

import (
	"crypto/x509"
	"fmt"
	"strings"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog"
	configv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
)

func (c *authOperator) handleRoute() (*routev1.Route, *corev1.Secret, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	route, err := c.route.Get(targetName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		route, err = c.route.Create(defaultRoute())
	}
	if err != nil {
		return nil, nil, err
	}
	if len(route.Spec.Host) == 0 {
		return nil, nil, fmt.Errorf("route has no host: %#v", route)
	}
	if err := isValidRoute(route); err != nil {
		klog.Infof("deleting invalid route: %#v", route)
		opts := &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &route.UID}}
		if err := c.route.Delete(route.Name, opts); err != nil && !errors.IsNotFound(err) {
			klog.Infof("failed to delete invalid route: %v", err)
		}
		return nil, nil, err
	}
	routerSecret, err := c.secrets.Secrets(targetNamespace).Get(routerCertsLocalName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	if len(routerSecret.Data) == 0 {
		return nil, nil, fmt.Errorf("router secret is empty: %#v", routerSecret)
	}
	return route, routerSecret, nil
}
func isValidRoute(route *routev1.Route) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expectedRoute := defaultRoute()
	expName := expectedRoute.Spec.To.Name
	expPort := expectedRoute.Spec.Port.TargetPort.IntValue()
	expTLSTermination := expectedRoute.Spec.TLS.Termination
	expInsecureEdgeTerminationPolicy := expectedRoute.Spec.TLS.InsecureEdgeTerminationPolicy
	if route.Spec.To.Name != expName {
		return fmt.Errorf("route targets a wrong service - needs %s: %#v", expName, route)
	}
	if route.Spec.Port.TargetPort.IntValue() != expPort {
		return fmt.Errorf("expected port '%d' for route: %#v", expPort, route)
	}
	if route.Spec.TLS == nil {
		return fmt.Errorf("TLS needs to be configured for route: %#v", route)
	}
	if route.Spec.TLS.Termination != expTLSTermination {
		return fmt.Errorf("route contains wrong TLS termination - '%s' is required: %#v", expTLSTermination, route)
	}
	if route.Spec.TLS.InsecureEdgeTerminationPolicy != expInsecureEdgeTerminationPolicy {
		return fmt.Errorf("route contains wrong insecure termination policy - '%s' is required: %#v", expInsecureEdgeTerminationPolicy, route)
	}
	return nil
}
func defaultRoute() *routev1.Route {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &routev1.Route{ObjectMeta: defaultMeta(), Spec: routev1.RouteSpec{To: routev1.RouteTargetReference{Kind: "Service", Name: targetName}, Port: &routev1.RoutePort{TargetPort: intstr.FromInt(containerPort)}, TLS: &routev1.TLSConfig{Termination: routev1.TLSTerminationPassthrough, InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect}}}
}
func routerSecretToSNI(routerSecret *corev1.Secret) []configv1.NamedCertificate {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var out []configv1.NamedCertificate
	for domain := range routerSecret.Data {
		out = append(out, configv1.NamedCertificate{Names: []string{"*." + domain}, CertInfo: configv1.CertInfo{CertFile: routerCertsLocalMount + "/" + domain, KeyFile: routerCertsLocalMount + "/" + domain}})
	}
	return out
}
func routerSecretToCA(route *routev1.Route, routerSecret *corev1.Secret) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		caData		[]byte
		longestDomain	string
	)
	for domain, certs := range routerSecret.Data {
		if strings.HasSuffix(route.Spec.Host, "."+domain) && len(domain) > len(longestDomain) {
			caData = certs
			longestDomain = domain
		}
	}
	if ok := x509.NewCertPool().AppendCertsFromPEM(caData); !ok {
		klog.Infof("using global CAs for %s, ingress domain=%s, cert data len=%d", route.Spec.Host, longestDomain, len(caData))
		return nil
	}
	return caData
}
