package commands

import (
	"testing"

	test "github.com/smartystreets/goconvey/convey"
)

func TestIssueCmd(t *testing.T) {
	test.Convey("test issue", t, func() {
		args := []string{"issue"}
		RootCmd.SetArgs(args)
		test.ShouldBeNil(RootCmd.Execute())
	})
}
