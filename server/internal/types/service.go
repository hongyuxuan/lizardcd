package types

import "github.com/hongyuxuan/lizardcd/agent/lizardagent"

type RpcAgent struct {
	Client        lizardagent.LizardAgent
	ServiceSource string
}
