package operator

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift/cluster-authentication-operator/pkg/boilerplate/controller"
)

type Option func(*operator)

func WithInformer(getter controller.InformerGetter, filter controller.Filter, opts ...controller.InformerOption) Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(o *operator) {
		o.opts = append(o.opts, controller.WithInformer(getter, controller.FilterFuncs{ParentFunc: func(obj v1.Object) (namespace, name string) {
			return o.name, o.name
		}, AddFunc: filter.Add, UpdateFunc: filter.Update, DeleteFunc: filter.Delete}, opts...))
	}
}
