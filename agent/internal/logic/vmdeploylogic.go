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
	"strconv"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/imroc/req/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type VmDeployLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVmDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VmDeployLogic {
	return &VmDeployLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
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
	}
	uid, _ := strconv.Atoi(deployUser.Uid)
	gid, _ := strconv.Atoi(deployUser.Gid)

	// save pre_comamnd to shell scripts
	if in.PreCommand != "" {
		prefile := fmt.Sprintf("%s/tmp-pre-command-%d.sh", in.DeployPath, time.Now().UnixMicro())
		if err = os.WriteFile(prefile, []byte(in.PreCommand), 0755); err != nil {
			l.Logger.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		// execute pre_command
		var cmd *exec.Cmd
		if currentUser.Name == "root" && in.DeployUser != "" { // only root chown, others not
			os.Chown(prefile, uid, gid)
			cmd = exec.Command("sudo", "-u", in.DeployUser, "-s", "/bin/bash", prefile)
		} else {
			cmd = exec.Command("/bin/bash", "-c", prefile)
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
		defer os.Remove(prefile)
		if err != nil {
			e := errorx.NewDefaultError("Failed to run pre_command: %v, output: %s", err, out.String())
			l.Logger.Error(e.Error())
			return nil, status.Error(codes.Internal, e.Error())
		}
		l.Logger.Infof("Successfully run pre_command \"%s\", output: %s", in.PreCommand, out.String())
	}

	// download package
	filename := filepath.Base(in.ArtifactUrl)
	if _, err = req.C().SetOutputDirectory(in.DeployPath).R().SetOutputFile(filename).SetHeaders(artifactHeaders).Get(in.ArtifactUrl); err != nil {
		l.Logger.Errorf("Failed to download package %s: %v", in.ArtifactUrl, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to download package %s: %v", in.ArtifactUrl, err))
	}

	// unarchive package
	err = utils.Unarchive(in.DeployPath+"/"+filename, in.DeployPath, uid, gid)
	defer os.Remove(in.DeployPath + "/" + filename)
	if err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	l.Logger.Infof("Unarchived package %s to %s", filename, in.DeployPath)

	// save start_comamnd to shell scripts
	shellfile := fmt.Sprintf("%s/tmp-start-command-%d.sh", in.DeployPath, time.Now().UnixMicro())
	if err = os.WriteFile(shellfile, []byte(in.StartCommand), 0755); err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// execute start_command
	var cmd *exec.Cmd
	if currentUser.Name == "root" && in.DeployUser != "" { // only root chown, others not
		os.Chown(shellfile, uid, gid)
		cmd = exec.Command("sudo", "-u", in.DeployUser, "-s", "/bin/bash", shellfile)
	} else {
		cmd = exec.Command("/bin/bash", "-c", shellfile)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	defer os.Remove(shellfile)
	if err != nil {
		e := errorx.NewDefaultError("Failed to run start_command: %v, output: %s", err, out.String())
		l.Logger.Error(e.Error())
		return nil, status.Error(codes.Internal, e.Error())
	}
	l.Logger.Infof("Successfully run start_command \"%s\", output: %s", in.StartCommand, out.String())
	resp = &agent.Response{
		Code: uint32(codes.OK),
		Data: out.Bytes(),
	}
	return
}
