package txretry

import (
	"context"
	"errors"
	"testing"
)

type fakeTx struct {
	commitErr    error
	committed    bool
	rolledBack   bool
	rollbackCall int
}

func (tx *fakeTx) Commit() error {
	tx.committed = true
	return tx.commitErr
}

func (tx *fakeTx) Rollback() error {
	tx.rollbackCall++
	tx.rolledBack = true
	return nil
}

func TestDoRetriesRetryableOperationError(t *testing.T) {
	t.Parallel()

	retryable := errors.New("serialization failure")
	var beginCalls int

	runner := Runner{
		MaxAttempts: 3,
		Begin: func(context.Context) (Tx, error) {
			beginCalls++
			return &fakeTx{}, nil
		},
		IsRetryable: func(err error) bool {
			return errors.Is(err, retryable)
		},
	}

	var attempts int
	err := runner.Do(context.Background(), func(context.Context, Tx) error {
		attempts++
		if attempts < 3 {
			return retryable
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if attempts != 3 || beginCalls != 3 {
		t.Fatalf("attempts = %d, beginCalls = %d, want 3", attempts, beginCalls)
	}
}

func TestDoStopsOnPermanentError(t *testing.T) {
	t.Parallel()

	permanent := errors.New("unique conflict")
	tx := &fakeTx{}

	runner := Runner{
		MaxAttempts: 3,
		Begin: func(context.Context) (Tx, error) {
			return tx, nil
		},
	}

	err := runner.Do(context.Background(), func(context.Context, Tx) error {
		return permanent
	})
	if !errors.Is(err, permanent) {
		t.Fatalf("Do() error = %v, want %v", err, permanent)
	}
	if !tx.rolledBack {
		t.Fatalf("Rollback() was not called")
	}
}

func TestDoRetriesRetryableCommitError(t *testing.T) {
	t.Parallel()

	retryable := errors.New("deadlock")
	txs := []*fakeTx{
		{commitErr: retryable},
		{},
	}
	var beginCalls int

	runner := Runner{
		MaxAttempts: 2,
		Begin: func(context.Context) (Tx, error) {
			tx := txs[beginCalls]
			beginCalls++
			return tx, nil
		},
		IsRetryable: func(err error) bool {
			return errors.Is(err, retryable)
		},
	}

	err := runner.Do(context.Background(), func(context.Context, Tx) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if beginCalls != 2 {
		t.Fatalf("beginCalls = %d, want 2", beginCalls)
	}
	if !txs[0].rolledBack {
		t.Fatalf("first tx should rollback after commit failure")
	}
}
