package solaris_telegraf_helpers

import (
	"fmt"
	"github.com/siebenmann/go-kstat"
	"log"
	"strconv"
	"strings"
)

// Return a list of IO kstats in the given class
func KstatIoClass(token *kstat.Token, class string) map[string]*kstat.IO {
	ret := make(map[string]*kstat.IO)

	for _, n := range token.All() {
		if n.Class != class {
			continue
		}

		stat, err := n.GetIO()

		if err != nil {
			log.Fatal("cannot get kstat")
		}

		key := fmt.Sprintf("%s:%s", n.Module, n.Name)

		ret[key] = stat
	}

	return ret
}

func KstatSingle(token *kstat.Token, statname string) interface{} {
	c := strings.Split(statname, ":")

	instance, _ := strconv.Atoi(c[1])
	stat, _ := token.GetNamed(c[0], instance, c[2], c[3])
	stat_type := fmt.Sprintf("%s", stat.Type)

	if stat_type == "uint32" || stat_type == "uint64" {
		return stat.UintVal
	}

	return stat.IntVal
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

func KstatClass(token *kstat.Token, class string) []*kstat.KStat {

	var ret []*kstat.KStat

	for _, stat := range token.All() {
		if stat.Class == class {
			ret = append(ret, stat)
		}
	}

	return ret
}

func KstatModule(token *kstat.Token, module string) []*kstat.KStat {

	var ret []*kstat.KStat

	for _, stat := range token.All() {
		if stat.Module == module {
			ret = append(ret, stat)
		}
	}

	return ret
}
