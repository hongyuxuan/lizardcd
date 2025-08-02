package logic

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunCommandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunCommandLogic {
	return &RunCommandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RunCommandLogic) RunCommand(in *agent.RunCommandRequest, stream agent.LizardAgent_RunCommandServer) (err error) {
	// get uid and gid
	runUser, _ := user.Current()
	if in.RunUser != "" {
		runUser, err = user.Lookup(in.RunUser)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("Cannot find user: %s", in.RunUser))
		}
		l.Logger.Infof("Using user: %s to run command", in.RunUser)
	}
	uid, _ := strconv.Atoi(runUser.Uid)
	gid, _ := strconv.Atoi(runUser.Gid)

	// save comamnd to shell scripts
	var commandname string
	if runtime.GOOS != "windows" {
		commandname = "command-*.sh"
	} else {
		if in.CommandType == "ps1" {
			commandname = "command-*.ps1"
		} else {
			commandname = "command-*.bat"
		}
	}
	var commandfile *os.File
	if commandfile, err = os.CreateTemp(os.TempDir(), commandname); err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	defer func() {
		if !in.ReserveTmp {
			os.Remove(commandfile.Name())
		}
	}()
	if in.Command != "" {
		_, err = commandfile.Write([]byte(in.Command))
	} else {
		_, err = commandfile.Write(in.ByteCommand)
	}
	if err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	l.Logger.Infof("Write command file to: %s", commandfile.Name())
	commandfile.Close()
	// execute command
	var cmd *exec.Cmd
	if in.RunUser != "" {
		if runtime.GOOS != "windows" {
			os.Chown(commandfile.Name(), uid, gid)
			cmd = exec.Command("sudo", "-u", in.RunUser, "-s", "/bin/sh", commandfile.Name())
		} else {
			if in.CommandType == "ps1" {
				cmd = exec.Command("powershell", "-File", commandfile.Name())
			} else {
				cmd = exec.Command("cmd", "/c", commandfile.Name())
			}
		}
	} else {
		if runtime.GOOS != "windows" {
			cmd = exec.Command("/bin/sh", commandfile.Name())
		} else {
			if in.CommandType == "ps1" {
				cmd = exec.Command("powershell", "-File", commandfile.Name())
			} else {
				cmd = exec.Command("cmd", "/c", commandfile.Name())
			}
		}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	if err = cmd.Start(); err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	go func() { // 处理 stdout
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if runtime.GOOS == "windows" {
				utf8Output, _, _ := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), scanner.Bytes())
				if err = stream.Send(&agent.YamlResponse{
					Code: uint32(codes.OK),
					Data: string(utf8Output),
				}); err != nil {
					l.Logger.Error(err)
				}
			} else {
				validOutput := scanner.Text()
				if !utf8.ValidString(validOutput) { // 非法utf8会造成对端主动关闭连接，报错transport is closing（gRPC 协议要求所有字符串字段必须是有效的 UTF-8 编码）
					validOutput = strings.ToValidUTF8(validOutput, "")
				}
				if err = stream.Send(&agent.YamlResponse{
					Code: uint32(codes.OK),
					Data: validOutput,
				}); err != nil {
					l.Logger.Error(err)
				}
			}
		}
	}()
	go func() { // 处理 stderr
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if runtime.GOOS == "windows" {
				utf8Output, _, _ := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), scanner.Bytes())
				if err = stream.Send(&agent.YamlResponse{
					Code: uint32(codes.OK),
					Data: string(utf8Output),
				}); err != nil {
					l.Logger.Error(err)
				}
			} else {
				validOutput := scanner.Text()
				if !utf8.ValidString(validOutput) { // 非法utf8会造成对端主动关闭连接，报错transport is closing（gRPC 协议要求所有字符串字段必须是有效的 UTF-8 编码）
					validOutput = strings.ToValidUTF8(validOutput, "")
				}
				if err = stream.Send(&agent.YamlResponse{
					Code: uint32(codes.OK),
					Data: validOutput,
				}); err != nil {
					l.Logger.Error(err)
				}
			}
		}
	}()
	// run finished
	if err = cmd.Wait(); err != nil {
		l.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	l.Logger.Infof("Run command finished.")
	return nil
}
