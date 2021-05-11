package solaris_telegraf_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Parallel()

	echoOutput := RunCmd("/bin/echo something")
	assert.Equal(t, "something", echoOutput)

	nosuch := RunCmd("/bin/no_such_command")
	assert.Equal(t, "", nosuch)
}
