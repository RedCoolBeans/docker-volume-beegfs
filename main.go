package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"

	"github.com/docker/go-plugins-helpers/volume"
	"os/user"
	"strconv"
)

var (
	// This is the path in beegfs-mounts.conf
	root    = flag.String("root", "/mnt/beegfs", "Base directory where volumes are created in the cluster")
	verbose = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)

	d := newBeeGFSDriver(*root)
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix("beegfs", gid))
}
