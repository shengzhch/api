package util

import (
	"api/log"
	"flag"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "Where to write CPU profile")

// Run starts up stuff at the beginning of a main function, and returns a
// function to defer until the function completes.  It should be used like this:
//
//   func main() {
//     defer util.Run()()
//     ... stuff ...
//   }

func Run() func() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("could not open cpu profile file %q", *cpuprofile)
		}
		_ = pprof.StartCPUProfile(f)
		return func() {
			pprof.StopCPUProfile()
			_ = f.Close()
		}
	}
	return func() {}
}
