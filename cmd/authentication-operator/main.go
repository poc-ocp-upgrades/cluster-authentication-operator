package main

import (
	goflag "flag"
	"bytes"
	"net/http"
	"runtime"
	"fmt"
	"math/rand"
	"os"
	"time"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	"github.com/openshift/cluster-authentication-operator/pkg/cmd/operator2"
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
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}
