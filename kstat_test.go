package solaris_telegraf_helpers

import (
	//"fmt"
	//sth "github.com/snltd/solaris-telegraf-helpers"
	"encoding/gob"
	"fmt"
	"github.com/siebenmann/go-kstat"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestKStatsInClass(t *testing.T) {
	allKStats = func(token *kstat.Token) []*kstat.KStat {
		return allKStatsFromFixtures()
	}

	var statNames []string

	for _, stat := range KStatsInClass(&kstat.Token{}, "errorq") {
		statNames = append(statNames, stat.Name)
	}

	assert.Equal(t, []string{"fm_ereport_queue", "gcpu_mca_queue"}, statNames)

	assert.Equal(
		t,
		[]*kstat.KStat(nil),
		KStatsInClass(&kstat.Token{}, "no_such_thing"),
	)
}

func TestKStatsInModule(t *testing.T) {
	allKStats = func(token *kstat.Token) []*kstat.KStat {
		return allKStatsFromFixtures()
	}

	for _, stat := range KStatsInModule(&kstat.Token{}, "cpu") {
		assert.Equal(t, "cpu", stat.Module)
	}

	assert.Equal(t, 12, len(KStatsInModule(&kstat.Token{}, "cpu")))

	assert.Equal(
		t,
		[]*kstat.KStat(nil),
		KStatsInModule(&kstat.Token{}, "no_such_thing"),
	)
}

func allKStatsFromFixtures() []*kstat.KStat {
	var kstatData []*kstat.KStat
	raw, err := os.Open("fixtures/all.kstat")

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not load serialized data from disk: %v\n", err))
	}

	dec := gob.NewDecoder(raw)
	err = dec.Decode(&kstatData)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load decode kstat data: %v\n", err)
		os.Exit(1)
	}

	return kstatData
}