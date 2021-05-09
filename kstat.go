package solaris_telegraf_helpers

import (
	"fmt"
	"github.com/siebenmann/go-kstat"
	"log"
	//"strconv"
	//"strings"
)

// NamedValue returns the useable value of the given named kstat. If said value is numeric, it is
// sent as a float64, which is what Telegraf expects as a value
func NamedValue(stat *kstat.Named) interface{} {
	switch stat.Type.String() {
	case "string", "char":
		return stat.StringVal
	case "int32", "int64":
		return float64(stat.IntVal)
	case "uint32", "uint64":
		return float64(stat.UintVal)
	default:
		log.Fatal(fmt.Sprintf("%s is of type %s", stat.Name, stat.Type))
	}

	return nil
}

// KStatIoClass returns a map of module:name => kstat for IO kstats
func KStatIoClass(token *kstat.Token, class string) map[string]*kstat.IO {
	ret := make(map[string]*kstat.IO)

	for _, n := range token.All() {
		if n.Class != class {
			continue
		}

		stat, err := n.GetIO()

		if err != nil {
			log.Fatal("cannot get kstat")
		}

		ret[fmt.Sprintf("%s:%s", n.Module, n.Name)] = stat
	}

	return ret
}

/*
// NamedStatValue returns the value of single, fully qualified, kstat. For instance,
// kstat cpu:0:sys:cpu_nsec_idle
func NamedStatValue(token *kstat.Token, statname string) interface{} {
	c := strings.Split(statname, ":")

	if len(c) != 4 {
		log.Fatal(fmt.Sprintf("'%s' is not a valid kstat", statname))
	}

	instance, _ := strconv.Atoi(c[1])
	stat, err := token.GetNamed(c[0], instance, c[2], c[3])

	if err != nil {
		log.Fatal(fmt.Sprintf("could not get value of %s", statname))
	}

	return NamedValue(stat)
}
*/

/*
// KstatBundle fetches a bundle of stats by name. For instance,
// kstat cpu:0:sys
func KstatBundle(token *kstat.Token, statname string) (*kstat.KStat, error) {
	c := strings.Split(statname, ":")

	if len(c) != 3 {
		log.Fatal(fmt.Sprintf("'%s' is not a valid kstat instance", statname))
	}

	instance, _ := strconv.Atoi(c[1])
	stat, err := token.GetNamed(c[0], instance, c[2], c[3])

	if err != nil {
		return *kstat.KStat{}, nil
	}

	return stat, nil
}

func KstatString(token *kstat.Token, statname string) string {
	c := strings.Split(statname, ":")

	instance, err := strconv.Atoi(c[1])

	if err != nil {
		return ""
	}

	v, err := token.GetNamed(c[0], instance, c[2], c[3])

	if err != nil {
		return ""
	}

	return v.StringVal
}
*/

// pulled out into a function to facilitate testing
var allKStats = func(token *kstat.Token) []*kstat.KStat {
	return token.All()
}

/*
// Wrapper for selected KStat information. Makes testing easier and code clearer IMO.
type kStatPath struct {
	Module   string
	Instance int
	Name     string
	Class    string
}
*/

// KStatsInClass returns a list of kstats in the given class
func KStatsInClass(token *kstat.Token, class string) []*kstat.KStat {
	var ret []*kstat.KStat

	for _, stat := range allKStats(token) {
		if stat.Class == class {
			ret = append(ret, stat)
		}
	}

	return ret
}

// KStatsInModule retuns a list of kstats in the given module. Asking for the 'cpu' module would
// give you something like
// cpu:0:intrstat
// cpu:0:sys
// cpu:0:vm
// cpu:1:intrstat
// ...
func KStatsInModule(token *kstat.Token, module string) []*kstat.KStat {
	var ret []*kstat.KStat

	for _, stat := range allKStats(token) {
		if stat.Module == module {
			ret = append(ret, stat)
		}
	}

	return ret
}
