package solaris_telegraf_helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZoneMapNames(t *testing.T) {
	zoneMap := ParseZones(zoneadmOutput)
	assert.ElementsMatch(t, []string{"global", "cube-media", "cube-ws"}, zoneMap.Names())
}

func TestZoneMapRunning(t *testing.T) {
	zoneMap := ParseZones(zoneadmOutput)
	assert.ElementsMatch(t, []string{"global", "cube-media"}, zoneMap.InState("running"))
	assert.ElementsMatch(t, []string{"cube-ws"}, zoneMap.InState("installed"))
	assert.ElementsMatch(t, []string{}, zoneMap.InState("configured"))
}

func TestParseZone(t *testing.T) {
	assert.Equal(
		t,
		zone{0, "global", "running", "/", "", "ipkg", "shared", 0},
		parseZone("0:global:running:/::ipkg:shared:0"),
	)

	assert.Equal(
		t,
		zone{
			42,
			"mz1",
			"running",
			"/zones/mz1",
			"c624d04f-d0d9-e1e6-822e-acebc78ec9ff",
			"lipkg",
			"excl",
			128,
		},
		parseZone("42:mz1:running:/zones/mz1:c624d04f-d0d9-e1e6-822e-acebc78ec9ff:lipkg:excl:128"),
	)
}

func TestParseZones(t *testing.T) {
	assert.Equal(
		t,
		ZoneMap{
			"global": zone{0, "global", "running", "/", "", "ipkg", "shared", 0},
			"cube-media": zone{
				42,
				"cube-media",
				"running",
				"/zones/cube-media",
				"c624d04f-d0d9-e1e6-822e-acebc78ec9ff",
				"lipkg",
				"excl",
				128,
			},
			"cube-ws": zone{
				44,
				"cube-ws",
				"installed",
				"/zones/cube-ws",
				"0f9c56f4-9810-6d45-f801-d34bf27cc13f",
				"pkgsrc",
				"excl",
				179,
			},
		},
		ParseZones(zoneadmOutput),
	)
}

var zoneadmOutput = `0:global:running:/::ipkg:shared:0
42:cube-media:running:/zones/cube-media:c624d04f-d0d9-e1e6-822e-acebc78ec9ff:lipkg:excl:128
44:cube-ws:installed:/zones/cube-ws:0f9c56f4-9810-6d45-f801-d34bf27cc13f:pkgsrc:excl:179`
