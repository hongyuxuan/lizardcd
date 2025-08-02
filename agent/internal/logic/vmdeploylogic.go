package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/imroc/req/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type VmDeployLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	vmService *commonsvc.VmService
}

func NewVmDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VmDeployLogic {
	return &VmDeployLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		vmService: commonsvc.NewVmService(ctx),
	}
}

// vm deploy
func (l *VmDeployLogic) VmDeploy(in *agent.VmDeployRequest) (resp *agent.Response, err error) {
	var artifactHeaders map[string]string
	if err = json.Unmarshal(in.ArtifactHeader, &artifactHeaders); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot unmarshal artifact_header into map[string]string: %v", err))
	}

	// get uid and gid
	currentUser, _ := user.Current()
	deployUser, _ := user.Current()
	if in.DeployUser != "" {
		deployUser, err = user.Lookup(in.DeployUser)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot find deploy user: %s", in.DeployUser))
		}
		l.Logger.Infof("Using user: %s to deploy", in.DeployUser)
	}
	uid, _ := strconv.Atoi(deployUser.Uid)
	gid, _ := strconv.Atoi(deployUser.Gid)

	// save pre_comamnd to shell scripts
	if in.PreCommand != "" {
		var prefile string
		if runtime.GOOS != "windows" {
			prefile = fmt.Sprintf("%s/tmp-pre-command-%d.sh", in.DeployPath, time.Now().UnixMicro())
		} else {
			if in.CommandType == "ps1" {
				prefile = fmt.Sprintf("%s/tmp-pre-command-%d.ps1", in.DeployPath, time.Now().UnixMicro())
			} else {
				prefile = fmt.Sprintf("%s/tmp-pre-command-%d.bat", in.DeployPath, time.Now().UnixMicro())
			}
		}
		if err = os.WriteFile(prefile, []byte(in.PreCommand), 0755); err != nil {
			l.Logger.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		l.Logger.Infof("Saved pre_command script to: %s", prefile)
		defer os.Remove(prefile)

		// execute pre_command
		var cmd *exec.Cmd
		if currentUser.Name == "root" && in.DeployUser != "" { // only root chown, others not
			os.Chown(prefile, uid, gid)
			cmd = exec.Command("sudo", "-u", in.DeployUser, "-s", "/bin/bash", prefile)
		} else {
			if runtime.GOOS != "windows" {
				cmd = exec.Command("/bin/bash", "-c", prefile)
			} else {
				if in.CommandType == "ps1" {
					cmd = exec.Command("powershell", "-File", prefile)
				} else {
					cmd = exec.Command("cmd", "/c", prefile)
				}
			}
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
		utf8Output := out.Bytes()
		if runtime.GOOS == "windows" {
			utf8Output, _, _ = transform.Bytes(simplifiedchinese.GBK.NewDecoder(), out.Bytes())
		}
		if err != nil {
			e := errorx.NewDefaultError("Failed to run pre_command: %v, output: %s", err, string(utf8Output))
			l.Logger.Error(e.Error())
			return nil, status.Error(codes.Internal, e.Error())
		}
		l.Logger.Infof("Successfully run pre_command, output: %s", string(utf8Output))
	}

	// download package
	filename := filepath.Base(in.ArtifactUrl)
	if _, err = req.C().SetOutputDirectory(in.DeployPath).R().SetOutputFile(filename).SetHeaders(artifactHeaders).Get(in.ArtifactUrl); err != nil {
		l.Logger.Errorf("Failed to download package %s: %v", in.ArtifactUrl, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to download package %s: %v", in.ArtifactUrl, err))
	}
	l.Logger.Infof("Successfully download package [%s] to %s", in.ArtifactUrl, in.DeployPath)

	// unarchive package
	if filepath.Ext(filename) != ".jar" {
		err = utils.Unarchive(in.DeployPath+"/"+filename, in.DeployPath, uid, gid)
		//defer os.Remove(in.DeployPath + "/" + filename)
		if err != nil {
			l.Logger.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		l.Logger.Infof("Unarchived package %s to %s", filename, in.DeployPath)
	}

	// save start_comamnd to shell scripts
	var commandfile string
	if runtime.GOOS != "windows" {
		commandfile = fmt.Sprintf("%s/tmp-start-command-%d.sh", in.DeployPath, time.Now().UnixMicro())
	} else {
		if in.CommandType == "ps1" {
			commandfile = fmt.Sprintf("%s/tmp-pre-command-%d.ps1", in.DeployPath, time.Now().UnixMicro())
		} else {
			commandfile = fmt.Sprintf("%s/tmp-pre-command-%d.bat", in.DeployPath, time.Now().UnixMicro())
		}
	}
	if err = os.WriteFile(commandfile, []byte(in.StartCommand), 0755); err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// execute start_command
	l.Logger.Infof("Saved start_command script to: %s", commandfile)
	var cmd *exec.Cmd
	if currentUser.Name == "root" && in.DeployUser != "" { // only root chown, others not
		os.Chown(commandfile, uid, gid)
		cmd = exec.Command("sudo", "-u", in.DeployUser, "-s", "/bin/bash", commandfile)
	} else {
		if runtime.GOOS != "windows" {
			cmd = exec.Command("/bin/bash", "-c", commandfile)
		} else {
			if in.CommandType == "ps1" {
				cmd = exec.Command("powershell", "-File", commandfile)
			} else {
				cmd = exec.Command("cmd", "/c", commandfile)
			}
		}
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	defer os.Remove(commandfile)
	utf8Output := out.Bytes()
	if runtime.GOOS == "windows" {
		utf8Output, _, _ = transform.Bytes(simplifiedchinese.GBK.NewDecoder(), out.Bytes())
	}
	if err != nil {
		e := errorx.NewDefaultError("Failed to run start_command: %v, output: %s", err, string(utf8Output))
		l.Logger.Error(e.Error())
		return nil, status.Error(codes.Internal, e.Error())
	}
	l.Logger.Infof("Successfully run start_command, output: %s", string(utf8Output))
	resp = &agent.Response{
		Code: uint32(codes.OK),
		Data: utf8Output,
	}
	return
}
