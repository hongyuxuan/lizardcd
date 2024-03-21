AppVersion=`cat ../VERSION`
GoVersion=`go env GOVERSION`
BuildTime=`date +"%Y-%m-%d"`
OsArch=`go env GOOS`/`go env GOARCH`
Author="hongyx2014@126.com"
Bin=`echo ${BINARY}`

LDFLAGS="-X 'github.com/hongyuxuan/lizardcd/cli/cmd.AppVersion=${AppVersion}' \
	-X 'github.com/hongyuxuan/lizardcd/cli/cmd.GoVersion=${GoVersion}' \
	-X 'github.com/hongyuxuan/lizardcd/cli/cmd.BuildTime=${BuildTime}' \
	-X 'github.com/hongyuxuan/lizardcd/cli/cmd.OsArch=${OsArch}' \
	-X 'github.com/hongyuxuan/lizardcd/cli/cmd.Author=${Author}'"

package:
	go build --ldflags ${LDFLAGS} -o bin/${Bin} main.go