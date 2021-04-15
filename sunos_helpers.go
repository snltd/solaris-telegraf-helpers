// Do Solaris/SmartOS stuff to help with my Telegraf collectors.
// Mostly wrapper methods to make life a little simpler.
package sunos_helpers

import (
	"fmt"
	"github.com/siebenmann/go-kstat"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Run a command, returning output as a string. Commands are
// specified as a simple string. You can't do pipes and stuff
// because of the way Go forks.
func RunCmd(cmd_str string) string {
	chunks := strings.SplitN(cmd_str, " ", 2)

	var c *exec.Cmd

	if len(chunks) == 2 {
		c = exec.Command(chunks[0], strings.Split(chunks[1], " ")...)
	} else {
		c = exec.Command(chunks[0])
	}

	p, _ := c.CombinedOutput()

	return strings.TrimSpace(string(p))
}

// Run a command via pfexec(1), returning its output as a string.
func RunCmdPfexec(cmd_str string) string {
	return RunCmd("pfexec " + cmd_str)
}

func Zone() string {
	return RunCmd("/bin/zonename")
}

func ZoneMap() map[int]string {
	zoneadm := ZoneadmRunning()

	ret := make(map[int]string, len(zoneadm))

	for _, line := range zoneadm {
		chunks := strings.Split(line, ":")
		key, _ := strconv.Atoi(chunks[0])
		ret[key] = chunks[1]
	}

	return ret
}

// Return the raw machine-parseable output from zoneadm, each zone
// in its own element.
func Zoneadm() []string {
	return strings.Split(RunCmd("/usr/sbin/zoneadm list -pc"), "\n")
}

func ZoneadmRunning() []string {
	return strings.Split(RunCmd("/usr/sbin/zoneadm list -p"), "\n")
}

// Return a list of zone names
func ZoneNames() []string {
	return strings.Split(RunCmd("/usr/sbin/zoneadm list"), "\n")
}

// Feed it a number with an ISO suffix and it will give you back the
// bytes in that number as a float.  'size' is a size such as '5G'
// or '0.5P'. (string)
func Bytify(size string) (float64, error) {
	return bytify_calc(size, 1024)
}

// As Bytify, but returns MiB, GiB etc.
func Bytify_t(size string) (float64, error) {
	return bytify_calc(size, 1000)
}

// Does the actual work for Bytify and Bytify_t.
func bytify_calc(size string, multiplier float64) (float64, error) {
	sizes := [8]string{"b", "K", "M", "G", "T", "P", "E", "Z"}

	if size == "-" {
		return 0, nil
	}

	r := regexp.MustCompile(`^\d+$`)

	if r.MatchString(size) {
		return strconv.ParseFloat(size, 64)
	}

	r = regexp.MustCompile(`^(-?[\d\.]+)(\w)$`)
	matches := r.FindAllStringSubmatch(size, -1)

	var exponent float64

	for i, v := range sizes {
		if v == matches[0][2] {
			exponent = float64(i)
			break
		}
	}

	base, _ := strconv.ParseFloat(matches[0][1], 64)

	return (base * (math.Pow(multiplier, exponent))), nil
}

func WeWant(want string, have []string) bool {

	if len(have) == 0 {
		return true
	}

	for _, thing := range have {
		if thing == want {
			return true
		}
	}

	return false
}

// Return a list of IO kstats in the given class
func KstatIoClass(token *kstat.Token, class string) map[string]*kstat.IO {
	ret := make(map[string]*kstat.IO)

	for _, n := range token.All() {
		if n.Class != class {
			continue
		}

		name := n.Name
		stat, err := n.GetIO()

		if err != nil {
			log.Fatal("cannot get kstat")
		}

		ret[name] = stat
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
