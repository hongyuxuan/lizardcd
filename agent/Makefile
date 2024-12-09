AppVersion=`cat ../VERSION`
GoVersion=`go env GOVERSION`
BuildTime=`date +'%F %T.%3N'`
OsArch=`go env GOOS`/`go env GOARCH`
Author="hongyx2014@126.com"
Bin=`echo ${BINARY}`

LDFLAGS="-X 'main.AppVersion=${AppVersion}' \
	-X 'main.GoVersion=${GoVersion}' \
	-X 'main.BuildTime=${BuildTime}' \
	-X 'main.OsArch=${OsArch}' \
	-X 'main.Author=${Author}'"

package:
	go build --ldflags ${LDFLAGS} -o bin/${Bin} agent.go