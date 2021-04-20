package solaris_telegraf_helpers

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCmd(t *testing.T) {
	echoOutput := RunCmd("/bin/echo something")
	assert.Equal(t, "something", echoOutput)

	nosuch := RunCmd("/bin/no_such_command")
	assert.Equal(t, "", nosuch)
}
