// Do Solaris/SmartOS stuff to help with my Telegraf collectors.
// Mostly wrapper methods to make life a little simpler.
package solaris_telegraf_helpers

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
