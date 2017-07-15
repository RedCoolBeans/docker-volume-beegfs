# docker-volume-beegfs

Docker Volume plugin to create persistent volumes in a [BeeGFS](http://www.beegfs.com/content/) cluster.

## Preconditions

- BeeGFS cluster has to be setup and running
- `beegfs-client` service needs to be running on the Docker host

## Installation

A pre-built binary as well as `rpm` and `deb` packages are available from the [releases](https://github.com/RedCoolBeans/docker-volume-beegfs/releases) page.

### RedHat/CentOS 7

An rpm can be built with:

    make rpm

Then install and start the service:

    yum localinstall docker-volume-beegfs-$VERSION.rpm
    systemctl start docker-volume-beegfs

### Debian 8

Debian packages are currently built on a RedHat system, but the `Makefile`
describes which packages to install on Debian when building from scratch.
Building the actual package can be done on a Debian system without Makefile modifications:

    make deb

Now you can install and start the service:

    dpkg -i docker-volume-beegfs_$VERSION.deb
    systemctl start docker-volume-beegfs

### From source code

The plugin uses [govendor](https://github.com/kardianos/govendor) to manage dependencies.

    go get -u github.com/kardianos/govendor

Restore dependencies:
    
    govendor sync

Build the plugin:

    go build

## Usage

First create a volume:

    docker volume create -d beegfs --name postgres-portroach

Then use the volume by passing the name (`postgres-1`):

    docker run -ti -v postgres-portroach:/var/lib/postgresql/data --volume-driver=beegfs -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres

Inspect the volume:

    docker volume inspect postgres-portroach

Remove the volume (note that this will _not_ remove the actual data):

    docker volume rm postgres-portroach

### Non-default mount points

By default BeeGFS uses `/mnt/beegfs` as the mount point (as configured in
`beegfs-mounts.conf`), and this plugin does too. For non-standard mount points
you can specify an alternate root when creating a new volume:

    docker volume create -d beegfs --name b3 -o root=/stor/b3

Other options are currently silently ignored.

## Roadmap

- No outstanding features/requests.

## License

MIT, please see the LICENSE file.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
