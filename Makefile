GIT_VER := $(shell git describe --tags)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)

clean:
	rm pkg/*

binary:
	gox -osarch="linux/amd64 darwin/amd64 windows/amd64 windows/386" -output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-X main.version ${GIT_VER} -X main.buildDate ${DATE}"
	cd pkg && find . -name "*${GIT_VER}*" -type f -exec zip {}.zip {} \;
