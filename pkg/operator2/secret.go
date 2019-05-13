package operator2

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func (c *authOperator) expectedSessionSecret() (*corev1.Secret, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret, err := c.secrets.Secrets(targetName).Get(sessionNameAndKey, metav1.GetOptions{})
	if err != nil || !isValidSessionSecret(secret) {
		klog.V(4).Infof("failed to get secret %s: %v", sessionNameAndKey, err)
		generatedSessionSecret, err := randomSessionSecret()
		if err != nil {
			return nil, err
		}
		return generatedSessionSecret, nil
	}
	return secret, nil
}
func isValidSessionSecret(secret *corev1.Secret) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var sessionSecretsBytes [][]byte
	for _, v := range secret.Data {
		sessionSecretsBytes = append(sessionSecretsBytes, v)
	}
	for _, ss := range sessionSecretsBytes {
		var sessionSecrets *legacyconfigv1.SessionSecrets
		err := json.Unmarshal(ss, &sessionSecrets)
		if err != nil {
			return false
		}
		for _, s := range sessionSecrets.Secrets {
			if len(s.Authentication) != 64 {
				return false
			}
			if len(s.Encryption) != 32 {
				return false
			}
		}
	}
	return true
}
func randomSessionSecret() (*corev1.Secret, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	skey, err := newSessionSecretsJSON()
	if err != nil {
		return nil, err
	}
	meta := defaultMeta()
	meta.Name = sessionNameAndKey
	return &corev1.Secret{ObjectMeta: meta, Data: map[string][]byte{sessionNameAndKey: skey}}, nil
}
func newSessionSecretsJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const (
		sha256KeyLenBytes = sha256.BlockSize
		aes256KeyLenBytes = 32
	)
	secrets := &legacyconfigv1.SessionSecrets{TypeMeta: metav1.TypeMeta{Kind: "SessionSecrets", APIVersion: "v1"}, Secrets: []legacyconfigv1.SessionSecret{{Authentication: randomString(sha256KeyLenBytes), Encryption: randomString(aes256KeyLenBytes)}}}
	secretsBytes, err := json.Marshal(secrets)
	if err != nil {
		return nil, fmt.Errorf("error marshalling the session secret: %v", err)
	}
	return secretsBytes, nil
}
func randomBytes(size int) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
func randomString(size int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b64size := base64.RawURLEncoding.DecodedLen(size) + 1
	return base64.RawURLEncoding.EncodeToString(randomBytes(b64size))[:size]
}
