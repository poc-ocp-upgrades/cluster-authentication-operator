package operator2

import (
	"fmt"
	"strings"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
)

func (c *authOperator) handleConfigSync(data *configSyncData) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	configMapClient := c.configMaps.ConfigMaps(targetNamespace)
	secretClient := c.secrets.Secrets(targetNamespace)
	configMaps, err := configMapClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing configMaps: %v", err)
	}
	secrets, err := secretClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing secrets: %v", err)
	}
	prefixConfigMapNames := sets.NewString()
	prefixSecretNames := sets.NewString()
	resourceVersionsAll := map[string]string{}
	for _, cm := range configMaps.Items {
		if strings.HasPrefix(cm.Name, userConfigPrefix) {
			prefixConfigMapNames.Insert(cm.Name)
			resourceVersionsAll[cm.Name] = cm.GetResourceVersion()
		}
	}
	for _, secret := range secrets.Items {
		if strings.HasPrefix(secret.Name, userConfigPrefix) {
			prefixSecretNames.Insert(secret.Name)
			resourceVersionsAll[secret.Name] = secret.GetResourceVersion()
		}
	}
	inUseConfigMapNames := sets.NewString()
	inUseSecretNames := sets.NewString()
	for dest, src := range data.idpConfigMaps {
		syncOrDie(c.resourceSyncer.SyncConfigMap, dest, src.src)
		inUseConfigMapNames.Insert(dest)
	}
	for dest, src := range data.idpSecrets {
		syncOrDie(c.resourceSyncer.SyncSecret, dest, src.src)
		inUseSecretNames.Insert(dest)
	}
	for dest, src := range data.tplSecrets {
		syncOrDie(c.resourceSyncer.SyncSecret, dest, src.src)
		inUseSecretNames.Insert(dest)
	}
	notInUseConfigMapNames := prefixConfigMapNames.Difference(inUseConfigMapNames)
	notInUseSecretNames := prefixSecretNames.Difference(inUseSecretNames)
	for dest := range notInUseConfigMapNames {
		syncOrDie(c.resourceSyncer.SyncConfigMap, dest, "")
	}
	for dest := range notInUseSecretNames {
		syncOrDie(c.resourceSyncer.SyncSecret, dest, "")
	}
	var resourceVersionsInUse []string
	for name := range inUseConfigMapNames {
		resourceVersionsInUse = append(resourceVersionsInUse, resourceVersionsAll[name])
	}
	for name := range inUseSecretNames {
		resourceVersionsInUse = append(resourceVersionsInUse, resourceVersionsAll[name])
	}
	return resourceVersionsInUse, nil
}

type configSyncData struct {
	idpConfigMaps	map[string]sourceData
	idpSecrets	map[string]sourceData
	tplSecrets	map[string]sourceData
}
type sourceData struct {
	src	string
	path	string
	volume	corev1.Volume
	mount	corev1.VolumeMount
}

func newSourceDataIDPSecret(index int, secretName configv1.SecretNameReference, field, key string) (string, sourceData) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dest := getIDPName(index, field)
	dirPath := getIDPPath(index, "secret", dest)
	vol, mount, path := secretVolume(dirPath, dest, key)
	ret := sourceData{src: secretName.Name, path: path, volume: vol, mount: mount}
	return dest, ret
}
func newSourceDataIDPConfigMap(index int, configMap configv1.ConfigMapNameReference, field, key string) (string, sourceData) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dest := getIDPName(index, field)
	dirPath := getIDPPath(index, "configmap", dest)
	vol, mount, path := configMapVolume(dirPath, dest, key)
	ret := sourceData{src: configMap.Name, path: path, volume: vol, mount: mount}
	return dest, ret
}
func newSourceDataTemplateSecret(secretRef configv1.SecretNameReference, field, key string) (string, sourceData) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dest := getTemplateName(field)
	dirPath := getTemplatePath("secret", dest)
	vol, mount, path := secretVolume(dirPath, dest, key)
	ret := sourceData{src: secretRef.Name, path: path, volume: vol, mount: mount}
	return dest, ret
}
func newConfigSyncData() configSyncData {
	_logClusterCodePath()
	defer _logClusterCodePath()
	idpConfigMaps := map[string]sourceData{}
	idpSecrets := map[string]sourceData{}
	tplSecrets := map[string]sourceData{}
	return configSyncData{idpConfigMaps: idpConfigMaps, idpSecrets: idpSecrets, tplSecrets: tplSecrets}
}
func (sd *configSyncData) addIDPSecret(index int, secretRef configv1.SecretNameReference, field, key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(secretRef.Name) == 0 {
		return ""
	}
	dest, data := newSourceDataIDPSecret(index, secretRef, field, key)
	sd.idpSecrets[dest] = data
	return data.path
}
func (sd *configSyncData) addIDPConfigMap(index int, configMapRef configv1.ConfigMapNameReference, field, key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(configMapRef.Name) == 0 {
		return ""
	}
	dest, data := newSourceDataIDPConfigMap(index, configMapRef, field, key)
	sd.idpConfigMaps[dest] = data
	return data.path
}
func (sd *configSyncData) addTemplateSecret(secretRef configv1.SecretNameReference, field, key string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(secretRef.Name) == 0 {
		return ""
	}
	dest, data := newSourceDataTemplateSecret(secretRef, field, key)
	sd.tplSecrets[dest] = data
	return data.path
}
func getIDPName(i int, field string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s%d-%s", userConfigPrefixIDP, i, field)
}
func getIDPPath(i int, resource, dest string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s/%d/%s/%s", userConfigPathPrefixIDP, i, resource, dest)
}
func getTemplateName(field string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return userConfigPrefixTemplate + field
}
func getTemplatePath(resource, dest string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s/%s/%s", userConfigPathPrefixTemplate, resource, dest)
}
func syncOrDie(syncFunc func(dest, src resourcesynccontroller.ResourceLocation) error, dest, src string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ns := userConfigNamespace
	if len(src) == 0 {
		ns = ""
	}
	if err := syncFunc(resourcesynccontroller.ResourceLocation{Namespace: targetNamespace, Name: dest}, resourcesynccontroller.ResourceLocation{Namespace: ns, Name: src}); err != nil {
		panic(err)
	}
}
func secretVolume(path, name, key string) (corev1.Volume, corev1.VolumeMount, string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data := volume{name: name, configmap: false, path: path, keys: []string{key}}
	vol, mount := data.split()
	return vol, mount, mount.MountPath + "/" + key
}
func configMapVolume(path, name, key string) (corev1.Volume, corev1.VolumeMount, string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data := volume{name: name, configmap: true, path: path, keys: []string{key}}
	vol, mount := data.split()
	return vol, mount, mount.MountPath + "/" + key
}
