all: deps compile

compile:
	go build

deps:
	go get

fmt:
	gofmt -s -w -l .

rpm-deps:
	yum install -y ruby ruby-devel rubygems rpm-build make go git
	gem install fpm

rpm: compile rpm-deps
	mkdir -p obj/redhat/usr/bin
	mkdir -p obj/redhat/lib/systemd/system/
	install -m 0755 docker-volume-beegfs obj/redhat/usr/bin
	install -m 0644 docker-volume-beegfs.service obj/redhat/lib/systemd/system
	fpm -C obj/redhat --vendor RedCoolBeans -m "info@redcoolbeans.com" -f \
		-s dir -t rpm -n docker-volume-beegfs \
		--after-install files/post-install-systemd --version 0.1.0 . && \
		rm -fr obj/redhat

# builds are done on RHEL, when building locally on Debian use the following:
# apt-get install -y ruby ruby-dev gcc golang git make
deb-deps:
	yum install -y ruby ruby-devel rubygems rpm-build make go git
	gem install fpm

deb: compile deb-deps
	mkdir -p obj/debian/usr/bin
	mkdir -p obj/debian/lib/systemd/system/
	install -m 0755 docker-volume-beegfs obj/debian/usr/bin
	install -m 0644 docker-volume-beegfs.service obj/debian/lib/systemd/system
	fpm -C obj/debian --vendor RedCoolBeans -m "info@redcoolbeans.com" -f \
		-s dir -t deb -n docker-volume-beegfs \
		--after-install files/post-install-systemd --version 0.1.0 . && \
		rm -fr obj/debian

clean:
	rm -fr obj *.deb *.rpm docker-volume-beegfs

.PHONY: clean rpm-deps deb-deps fmt deps compile
