package autoelect

import (
	"context"
)

type AutoElection interface {
	IsMaster() bool
	LoopInElect(ctx context.Context, errDuringLoopCh chan error) error
	StopElect()
}
