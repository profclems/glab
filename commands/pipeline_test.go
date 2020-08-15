package commands

import (
	test "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPipelineCmd(t *testing.T) {
	test.Convey("test pipeline cmd", t, func() {
		args := []string{"pipeline"}
		RootCmd.SetArgs(args)
		test.ShouldBeNil(RootCmd.Execute())
	})
	test.Convey("test pipeline alias cmd", t, func() {
		args := []string{"pipe"}
		RootCmd.SetArgs(args)
		test.ShouldBeNil(RootCmd.Execute())
	})
}
