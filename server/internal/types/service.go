package types

import (
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/zeromicro/go-zero/zrpc"
)

type RpcAgent struct {
	Client        lizardagent.LizardAgent
	ServiceSource string
	Cli           zrpc.Client
}
