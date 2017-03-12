//+build appengine

package cmd

import "os"

func signalNotify(interrupt chan<- os.Signal) {
	// Does not notify in the case of AppEngine.
}
