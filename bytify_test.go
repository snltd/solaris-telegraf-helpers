package solaris_telegraf_helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		result, _ := Bytify(table.in)
		assert.Equal(t, table.out, result)
	}
}

func TestBytifyI(t *testing.T) {
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
		result, _ := BytifyI(table.in)
		assert.Equal(t, table.out, result)
	}
}
