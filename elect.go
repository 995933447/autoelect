package autoelect

import (
	"context"
)

type AutoElection interface {
	LoopInElect(ctx context.Context)
	StopElect()
}
