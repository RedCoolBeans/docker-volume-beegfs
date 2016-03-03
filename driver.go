package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/davecgh/go-spew/spew"
)

type beegfsDriver struct {
	root    string
	volumes map[string]string
	m       *sync.Mutex
}

func newBeeGFSDriver(root string) beegfsDriver {
	d := beegfsDriver{
		root:    root,
		volumes: make(map[string]string),
		m:       &sync.Mutex{},
	}

	return d
}

func (b beegfsDriver) Create(r volume.Request) volume.Response {
	log.Infof("Create: %s, %v", r.Name, r.Options)

	b.m.Lock()
	defer b.m.Unlock()

	dest := volumeDir(b, r)

	if !isbeegfs(dest) {
		emsg := fmt.Sprintf("Cannot create volume %s as it's not on a BeeGFS filesystem", dest)
		log.Error(emsg)
		return volume.Response{Err: emsg}
	}

	if _, ok := b.volumes[r.Name]; ok {
		imsg := fmt.Sprintf("Cannot create volume %s, it already exists", dest)
		log.Info(imsg)
		return volume.Response{}
	}

	volumePath := volumeDir(b, r)

	if err := createDest(dest); err != nil {
		return volume.Response{Err: err.Error()}
	}

	b.volumes[r.Name] = volumePath

	if *verbose {
		spew.Dump(b.volumes)
	}

	return volume.Response{}
}

func (b beegfsDriver) Remove(r volume.Request) volume.Response {
	log.Infof("Remove: %s", r.Name)

	b.m.Lock()
	defer b.m.Unlock()

	if _, ok := b.volumes[r.Name]; ok {
		delete(b.volumes, r.Name)
	}

	return volume.Response{}
}

func (b beegfsDriver) Path(r volume.Request) volume.Response {
	log.Debugf("Path: %s", r.Name)

	if volumePath, ok := b.volumes[r.Name]; ok {
		return volume.Response{Mountpoint: volumePath}
	}

	return volume.Response{}
}

func (b beegfsDriver) Mount(r volume.Request) volume.Response {
	log.Infof("Mount: %s", r.Name)
	dest := volumeDir(b, r)

	if !isbeegfs(dest) {
		emsg := fmt.Sprintf("Cannot mount volume %s as it's not on a BeeGFS filesystem", dest)
		log.Error(emsg)
		return volume.Response{Err: emsg}
	}

	if volumePath, ok := b.volumes[r.Name]; ok {
		return volume.Response{Mountpoint: volumePath}
	}

	return volume.Response{}
}

func (b beegfsDriver) Unmount(r volume.Request) volume.Response {
	log.Infof("Unmount: %s", r.Name)
	return volume.Response{}
}

func (b beegfsDriver) Get(r volume.Request) volume.Response {
	log.Infof("Get: %s", r.Name)

	if volumePath, ok := b.volumes[r.Name]; ok {
		return volume.Response{
			Volume: &volume.Volume{
				Name:       r.Name,
				Mountpoint: volumePath,
			},
		}
	}

	return volume.Response{Err: fmt.Sprintf("volume %s unknown", r.Name)}
}

func (b beegfsDriver) List(r volume.Request) volume.Response {
	log.Infof("List %v", r)

	volumes := []*volume.Volume{}

	for name, path := range b.volumes {
		volumes = append(volumes, &volume.Volume{Name: name, Mountpoint: path})
	}

	return volume.Response{Volumes: volumes}
}

func volumeDir(b beegfsDriver, r volume.Request) string {
	// We should use a per volume type to keep track of their individual roots.
	// Then we can use r.Options["beegfsbase"]
	return filepath.Join(b.root, r.Name)
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
