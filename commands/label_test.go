package commands

import (
	"testing"

	test "github.com/smartystreets/goconvey/convey"
)

func TestLabelCmd(t *testing.T) {
	test.Convey("test label cmd", t, func() {
		args := []string{"label"}
		RootCmd.SetArgs(args)
		test.ShouldBeNil(RootCmd.Execute())
	})
}
