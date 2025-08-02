package logic

import (
	"context"
	"io"
	"sync"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

type RunCommandInPodLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunCommandInPodLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunCommandInPodLogic {
	return &RunCommandInPodLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RunCommandInPodLogic) RunCommandInPod(stream agent.LizardAgent_RunCommandInPodServer) error {
	msg, err := stream.Recv()
	if err != nil {
		if status.Code(err) == codes.Canceled {
			l.Logger.Errorf("流被取消: %v", err)
		} else {
			l.Logger.Error(err)
		}
		return status.Error(codes.Internal, err.Error())
	}
	req := l.svcCtx.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(msg.GetNamespace()).
		Name(msg.GetPodname()).
		SubResource("exec").
		Param("container", msg.GetContainerName()).
		VersionedParams(&corev1.PodExecOptions{
			Container: msg.GetContainerName(),
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(l.svcCtx.RestConfig, "POST", req.URL())
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(2)

	// 处理从客户端接收的输入
	go func() {
		defer wg.Done()
		defer stdinWriter.Close()
		for {
			msg, err = stream.Recv()
			l.Logger.Debugf("Received message: %+v", msg)
			if err != nil {
				if err != io.EOF {
					l.Logger.Errorf("Failed to receive from stream: %v", err)
				}
				return
			}
			if input := msg.GetCommand(); input != "" {
				input += "\n"
				if _, err := stdinWriter.Write([]byte(input)); err != nil {
					l.Logger.Errorf("Failed to write to stdin: %v", err)
					return
				}
			}
		}
	}()

	// // 处理输出流 (从 Kubernetes 到 gRPC)
	go func() {
		defer wg.Done()
		defer stdoutReader.Close()
		buf := make([]byte, 4096)
		for {
			n, err := stdoutReader.Read(buf)
			if err != nil {
				if err != io.EOF {
					l.Logger.Errorf("Failed to read from stdout: %v", err)
				}
				return
			}
			if err := stream.Send(&agent.Response{
				Code: uint32(codes.OK),
				Data: buf[:n],
			}); err != nil {
				l.Logger.Errorf("Failed to send to stream: %v", err)
				return
			}
		}
	}()

	// 执行远程命令
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  stdinReader,
		Stdout: stdoutWriter,
		Stderr: stdoutWriter,
		Tty:    true,
	})

	stdoutWriter.Close()
	stdinReader.Close()
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}
