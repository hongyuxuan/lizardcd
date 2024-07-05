package types

import (
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/zeromicro/go-zero/zrpc"
)

type RpcAgent struct {
	Client        lizardagent.LizardAgent
	ServiceSource string
	Cli           zrpc.Client
	Count         int
}

type HttpcheckResponse struct {
	Finished bool        `json:"finished"`
	Success  bool        `json:"success"`
	Message  interface{} `json:"message,omitempty"`
}
