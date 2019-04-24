package operator2

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"k8s.io/apimachinery/pkg/util/net"
	ktransport "k8s.io/client-go/transport"
)

func transportFor(caData, certData, keyData []byte) (http.RoundTripper, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	transport, err := transportForInner(caData, certData, keyData)
	if err != nil {
		return nil, err
	}
	return ktransport.DebugWrappers(transport), nil
}
func transportForInner(caData, certData, keyData []byte) (http.RoundTripper, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(caData) == 0 && len(certData) == 0 && len(keyData) == 0 {
		return http.DefaultTransport, nil
	}
	if (len(certData) == 0) != (len(keyData) == 0) {
		return nil, errors.New("cert and key data must be specified together")
	}
	transport := net.SetTransportDefaults(&http.Transport{TLSClientConfig: &tls.Config{}})
	if len(caData) != 0 {
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(caData); !ok {
			return nil, errors.New("error loading cert pool from ca data")
		}
		transport.TLSClientConfig.RootCAs = roots
	}
	if len(certData) != 0 {
		cert, err := tls.X509KeyPair(certData, keyData)
		if err != nil {
			return nil, errors.New("error loading x509 keypair from cert and key data")
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}
	return transport, nil
}
