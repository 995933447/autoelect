package gitdistribmu

import (
	"context"
	"errors"
	"fmt"
	"github.com/995933447/autoelect"
	"github.com/995933447/autoelect/util"
	"github.com/995933447/distribmu"
	"github.com/995933447/distribmu/factory"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//
// 基于分布式锁的主从选举
//
type AutoElection struct {
	mu       	distribmu.Mutex
	isMaster 	bool
	stopSignCh 	chan struct{}
}

var _ autoelect.AutoElection = (*AutoElection)(nil)

func (e *AutoElection) IsMaster() bool {
	return e.isMaster
}

func New(key string, ttl time.Duration, muConfType factory.MuType, muDriverConf interface{}) (*AutoElection, error) {
	if key == "" {
		return nil, errors.New("invalid key")
	}

	election := new(AutoElection)

	localIps, err := util.GetLocalIpsWithoutLoopback()
	if err != nil {
		return nil, err
	}

	if len(localIps) == 0 {
		return nil, errors.New("not found any local ip")
	}

	macAddrs, err := util.GetMacAddrs()
	if err != nil {
		return nil, err
	}

	if len(macAddrs) == 0 {
		return nil, errors.New("not found any mac address")
	}

	muConf := factory.NewMuConf(muConfType, key, ttl, fmt.Sprintf("ip=%s;mac=%s", localIps[0], macAddrs[0]), muDriverConf)
	election.mu = factory.MustNewMu(muConf)

	return election, nil
}

func (e *AutoElection) LoopInElect(ctx context.Context) {
	// 信号监听
	procSignChan := make(chan os.Signal, 1)
	go func() {
		for {
			sig := <-procSignChan
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				if e.isMaster {
					e.isMaster = false
					_ = e.mu.Unlock(ctx, false)
				}
				os.Exit(0)
			}
		}
	}()
	signal.Notify(procSignChan, syscall.SIGINT, syscall.SIGTERM)

	e.stopSignCh = make(chan struct{}, 1)
	defer func() {
		if e.isMaster {
			e.isMaster = false
			_ = e.mu.Unlock(ctx, false)
		}
	}()

	for {
		var isStop bool
		select {
		case _ = <-e.stopSignCh:
			isStop = true
		default:
		}

		if isStop {
			break
		}

		if e.isMaster {
			// 续期
			if time.Now().Add(8 * time.Second).After(e.mu.GetExpireTime()) {
				err := e.mu.RefreshTTL(ctx)
				if err != nil {
					if err != distribmu.ErrLockLost {
						err = e.mu.Unlock(ctx, false)
						if err != nil {
						}
					}
					// 刷新失败，当失去了 master 地位
					e.isMaster = false
				}
				continue
			}

			time.Sleep(time.Second)
			continue
		}

		locked, err := e.mu.LockWait(ctx, 10 * time.Second)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		if locked {
			e.isMaster = true
		}
	}
}

func (e *AutoElection) StopElect() {
	if e.stopSignCh != nil {
		e.stopSignCh <- struct{}{}
	}
}

