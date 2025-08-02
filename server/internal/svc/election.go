package svc

import (
	"context"
	"sync"
	"time"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type LeaderElection struct {
	client              *clientv3.Client
	session             *concurrency.Session
	election            *concurrency.Election
	isLeader            bool
	leaderMtx           sync.RWMutex
	instanceID          string
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	LeaderRetryInterval int
	LeaderRetryTimeout  int
}

func NewLeaderElection(cli *clientv3.Client, cfg config.Config) (*LeaderElection, error) {
	// 生成唯一实例ID
	instanceID := utils.GetLocalListener(cfg.Port)

	// 创建 session
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(cfg.Etcd.TTL))
	if err != nil {
		cli.Close()
		return nil, err
	}

	// 创建选举
	election := concurrency.NewElection(session, cfg.Etcd.LeaderPrefix)

	return &LeaderElection{
		client:              cli,
		session:             session,
		election:            election,
		instanceID:          instanceID,
		LeaderRetryInterval: cfg.Etcd.LeaderRetryInterval,
		LeaderRetryTimeout:  cfg.Etcd.LeaderRetryTimeout,
	}, nil
}

func (le *LeaderElection) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	le.cancel = cancel

	le.wg.Add(1)
	go func() {
		defer le.wg.Done()
		le.campaign(ctx)
	}()

	le.wg.Add(1)
	go func() {
		defer le.wg.Done()
		le.watchLeader(ctx)
	}()
}

func (le *LeaderElection) Stop() {
	if le.cancel != nil {
		le.cancel()
	}
	le.wg.Wait() // 等待所有goroutine结束

	if le.session != nil {
		le.session.Close()
	}
	if le.client != nil {
		le.client.Close()
	}
}

func (le *LeaderElection) IsLeader() bool {
	le.leaderMtx.RLock()
	defer le.leaderMtx.RUnlock()
	return le.isLeader
}

func (le *LeaderElection) campaign(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 尝试成为 leader，设置超时避免长期阻塞
			campaignCtx, cancel := context.WithTimeout(ctx, time.Duration(le.LeaderRetryTimeout)*time.Second)
			err := le.election.Campaign(campaignCtx, le.instanceID)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					logx.Debugf("Leader election timeout, already has a leader")
					time.Sleep(time.Duration(le.LeaderRetryInterval) * time.Second)
					continue
				}

				if ctx.Err() != nil {
					return
				}

				logx.Errorf("Leader election failed: %v", err)
				time.Sleep(time.Duration(le.LeaderRetryInterval) * time.Second)
				continue
			}

			// 成为 leader
			le.leaderMtx.Lock()
			le.isLeader = true
			le.leaderMtx.Unlock()
			logx.Infof("Instance=%s become leader", le.instanceID)

			// 保持 leader 状态
			select {
			case <-le.session.Done():
				le.leaderMtx.Lock()
				le.isLeader = false
				le.leaderMtx.Unlock()
				logx.Infof("Instance=%s is not leader", le.instanceID)
			case <-ctx.Done():
				// 主动退出时放弃领导权
				resignCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := le.election.Resign(resignCtx); err != nil {
					logx.Errorf("Resign leader failed: %v", err)
				}
				return
			}
		}
	}
}

func (le *LeaderElection) watchLeader(ctx context.Context) {
	ch := le.election.Observe(ctx)

	for {
		select {
		case resp, ok := <-ch:
			if !ok {
				logx.Error("leader watch channel close")
				return
			}
			if len(resp.Kvs) > 0 {
				logx.Infof("Now leader: %s", string(resp.Kvs[0].Value))
			} else {
				logx.Info("Now there is no leader")
			}
		case <-ctx.Done():
			return
		}
	}
}
