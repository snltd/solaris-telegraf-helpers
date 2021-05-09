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
	result, err := parseZone("0:global:running:/::ipkg:shared:0")
	assert.Nil(t, err)

	assert.Equal(
		t,
		zone{0, "global", "running", "/", "", "ipkg", "shared", 0},
		result,
	)

	result, err = parseZone(
		"42:mz1:running:/zones/mz1:c624d04f-d0d9-e1e6-822e-acebc78ec9ff:lipkg:excl:128",
	)

	assert.Nil(t, err)
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
		result,
	)

	result, err = parseZone("some:random:string")
	assert.Equal(t, zone{}, result)
	assert.Error(t, err)
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

func TestZoneByID(t *testing.T) {
	zoneMap := ParseZones(zoneadmOutput)
	assert.ElementsMatch(t, []string{"global", "cube-media", "cube-ws"}, zoneMap.Names())

	zoneData, err := zoneMap.ZoneByID(42)

	assert.Nil(t, err)
	assert.Equal(
		t,
		zone{
			ID:      42,
			Name:    "cube-media",
			Status:  "running",
			Path:    "/zones/cube-media",
			Uuid:    "c624d04f-d0d9-e1e6-822e-acebc78ec9ff",
			Brand:   "lipkg",
			IpType:  "excl",
			DebugID: 128,
		},
		zoneData)

	_, err = zoneMap.ZoneByID(101)
	assert.Error(t, err)
}

func TestParseZoneVnics(t *testing.T) {
	assert.Equal(
		t,
		ZoneVnicMap{
			"www_records0": vnic{
				Name:  "www_records0",
				Zone:  "cube-www-records",
				Link:  "rge0",
				Speed: 1000,
			},
		},
		ParseZoneVnics("www_records0:cube-www-records:rge0:1000"),
	)
}

func TestParseZoneVnic(t *testing.T) {
	assert.Equal(
		t,
		vnic{
			Name:  "www_records0",
			Zone:  "cube-www-records",
			Link:  "rge0",
			Speed: 1000,
		},
		parseZoneVnic("www_records0:cube-www-records:rge0:1000"),
	)
}

var zoneadmOutput = `0:global:running:/::ipkg:shared:0
42:cube-media:running:/zones/cube-media:c624d04f-d0d9-e1e6-822e-acebc78ec9ff:lipkg:excl:128
44:cube-ws:installed:/zones/cube-ws:0f9c56f4-9810-6d45-f801-d34bf27cc13f:pkgsrc:excl:179`
