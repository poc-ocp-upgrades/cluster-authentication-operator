package oscmd

import (
	"github.com/openshift/cluster-authentication-operator/tools/junitreport/pkg/builder"
	"github.com/openshift/cluster-authentication-operator/tools/junitreport/pkg/parser"
	"github.com/openshift/cluster-authentication-operator/tools/junitreport/pkg/parser/stack"
)

func NewParser(builder builder.TestSuitesBuilder, stream bool) parser.TestOutputParser {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return stack.NewParser(builder, newTestDataParser(), newTestSuiteDataParser(), stream)
}
