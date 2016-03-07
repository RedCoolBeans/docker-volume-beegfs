package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/go-plugins-helpers/volume"
)

// A single volume instance
type beegfsMount struct {
	name string
	path string
	root string
}

type beegfsDriver struct {
	mounts  map[string]*beegfsMount
	m       *sync.Mutex
}

func newBeeGFSDriver(root string) beegfsDriver {
	d := beegfsDriver{
		mounts:  make(map[string]*beegfsMount),
		m:       &sync.Mutex{},
	}

	return d
}

func (b beegfsDriver) Create(r volume.Request) volume.Response {
	var volumeRoot string

	log.Infof("Create: %s, %v", r.Name, r.Options)

	b.m.Lock()
	defer b.m.Unlock()

	// Handle options (unrecognized options are silently ignored):
	// root: directory to create new volumes (this should correspond with
	//       beegfs-mounts.conf).
	if optsRoot, ok := r.Options["root"]; ok {
		volumeRoot = optsRoot
	} else {
		// Assume the default root
		volumeRoot = *root
	}

	dest := filepath.Join(volumeRoot, r.Name)
	if !isbeegfs(dest) {
		emsg := fmt.Sprintf("Cannot create volume %s as it's not on a BeeGFS filesystem", dest)
		log.Error(emsg)
		return volume.Response{Err: emsg}
	}

	fmt.Printf("mounts: %d", len(b.mounts))
	if _, ok := b.mounts[r.Name]; ok {
		imsg := fmt.Sprintf("Cannot create volume %s, it already exists", dest)
		log.Info(imsg)
		return volume.Response{}
	}

	volumePath := filepath.Join(volumeRoot, r.Name)

	if err := createDest(dest); err != nil {
		return volume.Response{Err: err.Error()}
	}

	b.mounts[r.Name] = &beegfsMount {
		name: r.Name,
		path: volumePath,
		root: volumeRoot,
	}

	if *verbose {
		spew.Dump(b.mounts)
	}

	return volume.Response{}
}

func (b beegfsDriver) Remove(r volume.Request) volume.Response {
	log.Infof("Remove: %s", r.Name)

	b.m.Lock()
	defer b.m.Unlock()

	if _, ok := b.mounts[r.Name]; ok {
		delete(b.mounts, r.Name)
	}

	return volume.Response{}
}

func (b beegfsDriver) Path(r volume.Request) volume.Response {
	log.Debugf("Path: %s", r.Name)

	if _, ok := b.mounts[r.Name]; ok {
		return volume.Response{Mountpoint: b.mounts[r.Name].path}
	}

	return volume.Response{}
}

func (b beegfsDriver) Mount(r volume.Request) volume.Response {
	log.Infof("Mount: %s", r.Name)
	dest := filepath.Join(b.mounts[r.Name].root, r.Name)

	if !isbeegfs(dest) {
		emsg := fmt.Sprintf("Cannot mount volume %s as it's not on a BeeGFS filesystem", dest)
		log.Error(emsg)
		return volume.Response{Err: emsg}
	}

	if _, ok := b.mounts[r.Name]; ok {
		return volume.Response{Mountpoint: b.mounts[r.Name].path}
	}

	return volume.Response{}
}

func (b beegfsDriver) Unmount(r volume.Request) volume.Response {
	log.Infof("Unmount: %s", r.Name)
	return volume.Response{}
}

func (b beegfsDriver) Get(r volume.Request) volume.Response {
	log.Infof("Get: %s", r.Name)

	if v, ok := b.mounts[r.Name]; ok {
		return volume.Response{
			Volume: &volume.Volume{
				Name:       v.name,
				Mountpoint: v.path,
			},
		}
	}

	return volume.Response{Err: fmt.Sprintf("volume %s unknown", r.Name)}
}

func (b beegfsDriver) List(r volume.Request) volume.Response {
	log.Infof("List %v", r)

	volumes := []*volume.Volume{}

	for v := range b.mounts {
		volumes = append(volumes, &volume.Volume{Name: b.mounts[v].name, Mountpoint: b.mounts[v].path})
	}

	return volume.Response{Volumes: volumes}
}

// Check if the parent directory (where the volume will be created)
// is of type 'beegfs' using the BEEGFS_MAGIC value.
func isbeegfs(volumepath string) bool {
	log.Debugf("isbeegfs() for %s", volumepath)
	stat := syscall.Statfs_t{}
	err := syscall.Statfs(path.Dir(volumepath), &stat)
	if err != nil {
		log.Errorf("Could not determine filesystem type for %s: %s", volumepath, err)
		return false
	}

	log.Debugf("Type for %s: %d", volumepath, stat.Type)

	// BEEGFS_MAGIC 0x19830326
	return stat.Type == int64(428016422)
}

func createDest(dest string) error {
	fstat, err := os.Lstat(dest)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(dest, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if fstat != nil && !fstat.IsDir() {
		return fmt.Errorf("%v already exist and it's not a directory", dest)
	}

	return nil
}
