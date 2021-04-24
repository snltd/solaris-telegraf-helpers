package solaris_telegraf_helpers

import (
	"fmt"
	"strconv"
	"strings"
)

type zone struct {
	ID      int
	Name    string
	Status  string
	Path    string
	Uuid    string
	Brand   string
	IpType  string
	DebugID int
}

// ZoneMap maps the name of a zone to a zone struct containing all its zoneadm properties
type ZoneMap map[string]zone

// NewZoneMap creates a ZoneMap describing the current state of the system
func NewZoneMap() ZoneMap {
	raw := RunCmd("/usr/sbin/zoneadm list -cp")
	return ParseZones(raw)
}

// Names returns a list of zones in the map
func (z ZoneMap) Names() []string {
	zones := []string{}

	for zone := range z {
		zones = append(zones, zone)
	}

	return zones
}

// ZoneByID returns the zone with the given ID
func (z ZoneMap) ZoneByID(id int) (zone, error) {
	for _, zone := range z {
		if zone.ID == id {
			return zone, nil
		}
	}

	return zone{}, fmt.Errorf("no zone with ID %d", id)
}

// Names returns a list of zones in the map
func (z ZoneMap) InState(state string) []string {
	zones := []string{}

	for zone, data := range z {
		if data.Status == state {
			zones = append(zones, zone)
		}
	}

	return zones
}

// ZoneName returns the name of the current zone
func ZoneName() string {
	return RunCmd("/bin/zonename")
}

// ParseZones turns a chunk of raw `zoneadm list -p` output into a ZoneMap. It is public so
// Telegraf tests can use it
func ParseZones(raw string) ZoneMap {
	rawZones := strings.Split(raw, "\n")
	ret := ZoneMap{}

	for _, rawZone := range rawZones {
		zone := parseZone(rawZone)
		ret[zone.Name] = zone
	}

	return ret
}

// parseZone turns a line of raw `zoneadm list -p` output into a zone struct. The format of such a
// line is
// zoneid:zonename:state:zonepath:uuid:brand:ip-type:debugid
func parseZone(raw string) zone {
	chunks := strings.Split(raw, ":")
	zoneID, _ := strconv.Atoi(chunks[0])
	debugID, _ := strconv.Atoi(chunks[7])

	return zone{
		zoneID,
		chunks[1],
		chunks[2],
		chunks[3],
		chunks[4],
		chunks[5],
		chunks[6],
		debugID,
	}
}
