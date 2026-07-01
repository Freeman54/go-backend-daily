package fanoutquorum

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestReadReturnsWhenQuorumReached(t *testing.T) {
	t.Parallel()

	replicas := []Replica{
		func(ctx context.Context) (Response, error) {
			select {
			case <-time.After(10 * time.Millisecond):
				return Response{Replica: "a", Value: "v1"}, nil
			case <-ctx.Done():
				return Response{}, ctx.Err()
			}
		},
		func(ctx context.Context) (Response, error) {
			select {
			case <-time.After(20 * time.Millisecond):
				return Response{Replica: "b", Value: "v1"}, nil
			case <-ctx.Done():
				return Response{}, ctx.Err()
			}
		},
		func(ctx context.Context) (Response, error) {
			select {
			case <-time.After(time.Second):
				return Response{Replica: "slow", Value: "old"}, nil
			case <-ctx.Done():
				return Response{}, ctx.Err()
			}
		},
	}

	start := time.Now()
	got, err := Read(context.Background(), 2, replicas)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(Read()) = %d, want 2", len(got))
	}
	if time.Since(start) >= 200*time.Millisecond {
		t.Fatalf("Read() did not return after quorum")
	}
}

func TestReadFailsWhenQuorumBecomesImpossible(t *testing.T) {
	t.Parallel()

	replicas := []Replica{
		func(context.Context) (Response, error) { return Response{}, errors.New("boom") },
		func(context.Context) (Response, error) { return Response{Replica: "b", Value: "v1"}, nil },
		func(context.Context) (Response, error) { return Response{}, errors.New("boom") },
	}

	_, err := Read(context.Background(), 2, replicas)
	if !errors.Is(err, ErrQuorumUnavailable) {
		t.Fatalf("Read() error = %v, want ErrQuorumUnavailable", err)
	}
}

func TestReadRejectsInvalidQuorum(t *testing.T) {
	t.Parallel()

	_, err := Read(context.Background(), 2, []Replica{
		func(context.Context) (Response, error) { return Response{}, nil },
	})
	if !errors.Is(err, ErrInvalidQuorum) {
		t.Fatalf("Read() error = %v, want ErrInvalidQuorum", err)
	}
}
