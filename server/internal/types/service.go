package types

import (
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/robfig/cron/v3"
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

type WsMessage struct {
	MessageType string                 `json:"message_type"`
	MessageTo   string                 `json:"message_to,omitempty"`
	MessageData map[string]interface{} `json:"message_data"`
}

type SyncStatus struct {
	AppName        string                            `json:"app_name"`
	Failed         bool                              `json:"failed,omitempty"`
	Error          error                             `json:"message,omitempty"`
	Status         string                            `json:"status"`
	Reason         string                            `json:"reason"`
	LocalCommitId  string                            `json:"local_commitId"`
	RemoteCommitId string                            `json:"remote_commitId"`
	Comment        string                            `json:"comment"`
	Author         string                            `json:"author"`
	Resource       []commontypes.ApplicationResource `json:"resource,omitempty"`
}

type CronData struct {
	CronId    cron.EntryID
	HostPort  string
	Scheduler string
}
