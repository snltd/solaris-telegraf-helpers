package sunos_telegraf_helpers_test

import (
	sh "github.com/snltd/sunos_telegraf_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRunCmd(t *testing.T) {
	echoOutput := sh.RunCmd("/bin/echo something")
	assert.Equal(t, "something", echoOutput)

	nosuch := sh.RunCmd("/bin/no_such_command")
	assert.Equal(t, "", nosuch)
}

func TestBytify(t *testing.T) {
	tables := []struct {
		in  string
		out float64
	}{
		{"-", 0},
		{"5", 5},
		{"15b", 15},
		{"2K", 2048},
		{"2.5K", 2560},
		{"800M", 838860800},
		{"6.12G", 6571299962.88},
		{"0.5T", 549755813888},
	}

	for _, table := range tables {
		result, _ := sh.Bytify(table.in)
		assert.Equal(t, table.out, result)
	}
}

func TestBytify_t(t *testing.T) {
	tables := []struct {
		in  string
		out float64
	}{
		{"-", 0},
		{"5", 5},
		{"15b", 15},
		{"2K", 2000},
		{"2.5K", 2500},
		{"800M", 800000000},
		{"-6.12G", -6120000000},
		{"0.5T", 500000000000},
	}

	for _, table := range tables {
		result, _ := sh.Bytify_t(table.in)
		assert.Equal(t, table.out, result)
	}
}
