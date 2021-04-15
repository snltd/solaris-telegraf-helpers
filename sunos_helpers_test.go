package sunos_helpers_test

import (
	"testing"
	"github.com/snltd/sunos_helpers"
)

func TestBytify(t *testing.T) {
	tables := []struct {
		in string
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
		calc, _ := sunos_helpers.Bytify(table.in)

		if calc != table.out {
			t.Errorf("Expected %d, got %d", table.out, calc)
		}
	}
}

func TestBytify_t(t *testing.T) {
	tables := []struct {
		in string
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
		calc, _ := sunos_helpers.Bytify_t(table.in)

		if calc != table.out {
			t.Errorf("Expected %d, got %d", table.out, calc)
		}
	}
}
