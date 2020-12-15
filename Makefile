commit := ${shell git rev-parse HEAD}
branch := ${shell git rev-parse --abbrev-ref HEAD}
# tag := ${shell git describe --tag}
appkit_commit := ${shell git --git-dir=${GOPATH}/src/github.com/paypermint/appkit/.git rev-parse HEAD}
mskit_commit := ${shell git --git-dir=${GOPATH}/src/github.com/paypermint/mskit/.git rev-parse HEAD}

bridge-app-svc: main.go
	go build -ldflags "-X github.com/paypermint/appkit.Commit=${commit} -X github.com/paypermint/appkit.Branch=$(branch)  -X github.com/paypermint/appkit.AppkitCommit=${appkit_commit} -X github.com/paypermint/appkit.MskitCommit=${mskit_commit}"

bridge-app-svc.ubuntu: main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -X github.com/paypermint/appkit.Commit=${commit} -X github.com/paypermint/appkit.Branch=$(branch) -X github.com/paypermint/appkit.AppkitCommit=${appkit_commit} -X github.com/paypermint/appkit.MskitCommit=${mskit_commit} ' -a -tags netgo .

image.slim: clean bridge-app-svc.ubuntu
	docker build -f Dockerfile.slim -t paypermint/bridge-app-svc:slim-latest .

clean:
	go clean
