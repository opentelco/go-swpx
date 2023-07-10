package workflows

import (
	"testing"

	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

var testAct = activities.Activities{}
var testDevAct = device.Activities{}

type unitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (t *unitTestSuite) SetupTest() {
	t.env = t.NewTestWorkflowEnvironment()
}

func (t *unitTestSuite) AfterTest(suiteName, testName string) {
	t.env.AssertExpectations(t.T())
}

// Run all the tests in the suite
func Test_UnitTestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}
