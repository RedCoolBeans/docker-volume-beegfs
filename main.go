package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "os"
	// "path/filepath"

	"github.com/calavera/dkvolume"
)

var (
	// This is the path in beegfs-mounts.conf
	root    = flag.String("root", "/mnt/beegfs", "Base directory where volumes are created in the cluster")
	verbose = flag.Bool("verbose", true, "Enable verbose logging")
)

func main() {
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	d := newBeeGFSDriver(*root)
	h := dkvolume.NewHandler(d)
	fmt.Println(h.ServeUnix("root", "beegfs"))
}
