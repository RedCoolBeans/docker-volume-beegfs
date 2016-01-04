# docker-volume-beegfs

Docker Volume plugin to create persistent volumes in a [BeeGFS](http://www.beegfs.com/content/) cluster.

## Preconditions

- BeeGFS cluster has to be setup and running
- `beegfs-client` service needs to be running on the Docker host

## Installation

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

### Others

Build the plugin:

    go build

## Usage

First create a volume (no additional options are supported yet):

    docker volume create -d beegfs --name postgres-portroach

Then use the volume by passing the name (`postgres-1`):

    docker run -ti -v postgres-portroach:/var/lib/postgresql/data --volume-driver=beegfs -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres

Inspect the volume:

    docker volume inspect postgres-portroach

Remove the volume (note that this will _not_ remove the actual data):

    docker volume rm postgres-portroach

## Caveats

- Currently the plugin assumes the BeeGFS share is mounted on `/mnt/beegfs`

## Roadmap

- Support options passed at `docker volume create`

## License

MIT, please see the LICENSE file.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
