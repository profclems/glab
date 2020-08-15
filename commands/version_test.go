package commands

import (
	"testing"

	test "github.com/smartystreets/goconvey/convey"
)

func TestVersionCmd(t *testing.T) {
	test.Convey("test version cmd", t, func() {
		args := []string{"version"}
		RootCmd.SetArgs(args)
		test.ShouldBeNil(RootCmd.Execute())
	})
}