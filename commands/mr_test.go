package commands

import (
	"testing"

	test "github.com/smartystreets/goconvey/convey"
)

func TestMrCmd(t *testing.T) {
	test.Convey("test mr", t, func() {
		args := []string{"mr"}
		RootCmd.SetArgs(args)
		RootCmd.Execute()
	})
}
