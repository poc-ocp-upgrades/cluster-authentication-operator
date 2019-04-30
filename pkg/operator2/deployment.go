package operator2

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	operatorv1 "github.com/openshift/api/operator/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
)

func defaultDeployment(operatorConfig *operatorv1.Authentication, syncData *configSyncData, routerSecret *corev1.Secret, operatorDeployment *appsv1.Deployment, resourceVersions ...string) *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	replicas := int32(2)
	gracePeriod := int64(30)
	var (
		volumes	[]corev1.Volume
		mounts	[]corev1.VolumeMount
	)
	for _, data := range []volume{{name: sessionNameAndKey, configmap: false, path: sessionMount, keys: []string{sessionNameAndKey}}, {name: cliConfigNameAndKey, configmap: true, path: cliConfigMount, keys: []string{cliConfigNameAndKey}}, {name: servingCertName, configmap: false, path: servingCertMount, keys: []string{corev1.TLSCertKey, corev1.TLSPrivateKeyKey}}, {name: serviceCAName, configmap: true, path: serviceCAMount, keys: []string{serviceCAKey}}, {name: routerCertsLocalName, configmap: false, path: routerCertsLocalMount, keys: sets.StringKeySet(routerSecret.Data).List()}} {
		v, m := data.split()
		volumes = append(volumes, v)
		mounts = append(mounts, m)
	}
	volumes, mounts = toVolumesAndMounts(syncData.idpConfigMaps, volumes, mounts)
	volumes, mounts = toVolumesAndMounts(syncData.idpSecrets, volumes, mounts)
	volumes, mounts = toVolumesAndMounts(syncData.tplSecrets, volumes, mounts)
	sort.Strings(resourceVersions)
	rvs := strings.Join(resourceVersions, ",")
	rvsHash := sha512.Sum512([]byte(rvs))
	rvsHashStr := base64.RawURLEncoding.EncodeToString(rvsHash[:])
	meta := defaultMeta()
	meta.Annotations[deploymentVersionHashKey] = rvsHashStr
	deployment := &appsv1.Deployment{ObjectMeta: meta, Spec: appsv1.DeploymentSpec{Replicas: &replicas, Selector: &metav1.LabelSelector{MatchLabels: defaultLabels()}, Template: corev1.PodTemplateSpec{ObjectMeta: meta, Spec: corev1.PodSpec{ServiceAccountName: targetName, RestartPolicy: corev1.RestartPolicyAlways, SchedulerName: corev1.DefaultSchedulerName, TerminationGracePeriodSeconds: &gracePeriod, SecurityContext: &corev1.PodSecurityContext{}, Containers: []corev1.Container{{Image: oauthserverImage, ImagePullPolicy: getImagePullPolicy(operatorDeployment), Name: targetName, Command: []string{"hypershift", "openshift-osinserver", fmt.Sprintf("--config=%s", cliConfigPath), fmt.Sprintf("--v=%d", getLogLevel(operatorConfig.Spec.LogLevel))}, Ports: []corev1.ContainerPort{{Name: "https", Protocol: corev1.ProtocolTCP, ContainerPort: containerPort}}, VolumeMounts: mounts, ReadinessProbe: defaultProbe(), LivenessProbe: livenessProbe(), TerminationMessagePath: "/dev/termination-log", TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"), Resources: corev1.ResourceRequirements{Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceCPU: resource.MustParse("10m"), corev1.ResourceMemory: resource.MustParse("50Mi")}}}}, NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""}, PriorityClassName: "system-cluster-critical", Affinity: &corev1.Affinity{PodAntiAffinity: &corev1.PodAntiAffinity{PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{{Weight: 100, PodAffinityTerm: corev1.PodAffinityTerm{LabelSelector: &metav1.LabelSelector{MatchLabels: defaultLabels()}, TopologyKey: "kubernetes.io/hostname"}}}}}, Tolerations: []corev1.Toleration{{Operator: corev1.TolerationOpExists}}, Volumes: volumes}}}}
	return deployment
}
func defaultProbe() *corev1.Probe {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &corev1.Probe{Handler: corev1.Handler{HTTPGet: &corev1.HTTPGetAction{Path: "/healthz", Port: intstr.FromInt(containerPort), Scheme: corev1.URISchemeHTTPS}}, TimeoutSeconds: 1, PeriodSeconds: 10, SuccessThreshold: 1, FailureThreshold: 3}
}
func livenessProbe() *corev1.Probe {
	_logClusterCodePath()
	defer _logClusterCodePath()
	probe := defaultProbe()
	probe.InitialDelaySeconds = 30
	return probe
}
func toVolumesAndMounts(data map[string]sourceData, volumes []corev1.Volume, mounts []corev1.VolumeMount) ([]corev1.Volume, []corev1.VolumeMount) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	names := sets.StringKeySet(data).List()
	for _, name := range names {
		volumes = append(volumes, data[name].volume)
		mounts = append(mounts, data[name].mount)
	}
	return volumes, mounts
}
func getLogLevel(logLevel operatorv1.LogLevel) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch logLevel {
	case operatorv1.Normal, "":
		return 2
	case operatorv1.Debug:
		return 4
	case operatorv1.Trace:
		return 6
	case operatorv1.TraceAll:
		return 100
	default:
		return 0
	}
}
func getImagePullPolicy(operatorDeployment *appsv1.Deployment) corev1.PullPolicy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	containers := operatorDeployment.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		return corev1.PullIfNotPresent
	}
	return containers[0].ImagePullPolicy
}

type volume struct {
	name		string
	configmap	bool
	path		string
	keys		[]string
}

func (v *volume) split() (corev1.Volume, corev1.VolumeMount) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vol := corev1.Volume{Name: v.name}
	var items []corev1.KeyToPath
	for _, key := range v.keys {
		items = append(items, corev1.KeyToPath{Key: key, Path: key})
	}
	if v.configmap {
		vol.ConfigMap = &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: v.name}, Items: items}
	} else {
		vol.Secret = &corev1.SecretVolumeSource{SecretName: v.name, Items: items}
	}
	return vol, corev1.VolumeMount{Name: v.name, ReadOnly: true, MountPath: v.path}
}
