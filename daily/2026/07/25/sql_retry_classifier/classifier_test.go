package sql_retry_classifier

import (
	"errors"
	"testing"
)

type codedError string

func (e codedError) Error() string    { return string(e) }
func (e codedError) SQLState() string { return string(e) }

func TestIsRetryable_ClassifiesTransientSQLStates(t *testing.T) {
	for _, state := range []string{"40001", "40P01", "55P03"} {
		if !IsRetryable(codedError(state)) {
			t.Errorf("SQLSTATE %s should be retryable", state)
		}
	}
	if IsRetryable(codedError("23505")) {
		t.Fatal("唯一约束冲突不应重试")
	}
}

func TestIsRetryable_HandlesWrappedAndUnknownErrors(t *testing.T) {
	if !IsRetryable(errors.Join(errors.New("query failed"), codedError("40001"))) {
		t.Fatal("包装后的序列化失败应可重试")
	}
	if IsRetryable(errors.New("timeout text is not SQLSTATE")) {
		t.Fatal("未知错误不应重试")
	}
}
