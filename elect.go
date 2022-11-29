package autoelect

import (
	"context"
)

type AutoElection interface {
	LoopInAuthElect(ctx context.Context)
	StopElect()
}
