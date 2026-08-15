package sqlstatementtimeoutbudget

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrNoBudget = errors.New("no statement timeout budget")

// Budget 把请求剩余时间转换为数据库语句预算，并为提交、回滚和响应编码保留时间。
func Budget(ctx context.Context, now time.Time, maximum, reserve time.Duration) (time.Duration, error) {
	if maximum <= 0 || reserve < 0 {
		return 0, ErrNoBudget
	}
	budget := maximum
	if deadline, ok := ctx.Deadline(); ok {
		remaining := deadline.Sub(now) - reserve
		if remaining < budget {
			budget = remaining
		}
	}
	if budget < time.Millisecond {
		return 0, ErrNoBudget
	}
	// PostgreSQL 的毫秒配置会截断小数，主动向下取整，避免超过上游 deadline。
	return budget.Truncate(time.Millisecond), nil
}

// SetLocalSQL 生成仅在当前 PostgreSQL 事务生效的参数化设置语句。
func SetLocalSQL(budget time.Duration) (string, []any, error) {
	if budget < time.Millisecond {
		return "", nil, ErrNoBudget
	}
	milliseconds := budget.Milliseconds()
	return "SELECT set_config('statement_timeout', $1, true)", []any{fmt.Sprintf("%dms", milliseconds)}, nil
}
