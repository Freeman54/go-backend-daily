package fanoutquorum

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrInvalidQuorum     = errors.New("invalid quorum")
	ErrQuorumUnavailable = errors.New("quorum unavailable")
)

type Response struct {
	Replica string
	Value   string
}

type Replica func(context.Context) (Response, error)

type outcome struct {
	resp Response
	err  error
}

func Read(ctx context.Context, need int, replicas []Replica) ([]Response, error) {
	if need <= 0 || need > len(replicas) {
		return nil, ErrInvalidQuorum
	}

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan outcome, len(replicas))
	var wg sync.WaitGroup
	for _, replica := range replicas {
		wg.Add(1)
		go func(call Replica) {
			defer wg.Done()
			resp, err := call(childCtx)
			results <- outcome{resp: resp, err: err}
		}(replica)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	successes := make([]Response, 0, need)
	received := 0
	for result := range results {
		received++
		if result.err == nil {
			successes = append(successes, result.resp)
			if len(successes) == need {
				cancel()
				return successes, nil
			}
		}

		remaining := len(replicas) - received
		if len(successes)+remaining < need {
			cancel()
			return nil, ErrQuorumUnavailable
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, ErrQuorumUnavailable
}
