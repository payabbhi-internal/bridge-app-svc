commit := ${shell git rev-parse HEAD}
branch := ${shell git rev-parse --abbrev-ref HEAD}
tag := ${shell git describe --tag}
mskit_commit := ${shell git --git-dir=${GOPATH}/src/github.com/paypermint/mskit/.git rev-parse HEAD}

bridge-app-svc: main.go
	go build -ldflags "-X github.com/paypermint/mskit.Commit=${commit} -X github.com/paypermint/mskit.Branch=${branch} -X github.com/paypermint/mskit.Tag=${tag} -X github.com/paypermint/mskit.MskitCommit=${mskit_commit}"

bridge-app-svc.ubuntu: main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -X github.com/paypermint/mskit.Commit=${commit} -X github.com/paypermint/mskit.Branch=${branch} -X github.com/paypermint/mskit.Tag=${tag} -X github.com/paypermint/mskit.MskitCommit=${mskit_commit}' -a -tags netgo .

build:
	chmod +x run.sh
	chmod +x bridge-app-svc
	tar -czvf bridge-app-svc.tar.gz run.sh bridge-app-svc

local: clean bridge-app-svc build

release: clean bridge-app-svc.ubuntu build

image.slim: clean bridge-app-svc.ubuntu
	docker build -f Dockerfile.slim -t paypermint/bridge-app-svc:slim-latest .

clean:
	rm -rf bridge-app-svc.tar.gz || true
	go clean
