package refreshbackoff

import (
	"fmt"
	"time"
)

// Policy 控制缓存刷新失败后的指数退避。
type Policy struct {
	Base time.Duration
	Max  time.Duration
}

// Next 返回下一次允许刷新的时间。attempt 从 1 开始。
func (p Policy) Next(now time.Time, attempt int) (time.Time, error) {
	if p.Base <= 0 || p.Max < p.Base {
		return time.Time{}, fmt.Errorf("退避参数无效")
	}
	if attempt < 1 {
		return time.Time{}, fmt.Errorf("attempt 必须至少为 1")
	}

	delay := p.Base
	for i := 1; i < attempt && delay < p.Max; i++ {
		if delay > p.Max/2 {
			delay = p.Max
			break
		}
		delay *= 2
	}
	if delay > p.Max {
		delay = p.Max
	}
	return now.Add(delay), nil
}

// Ready 判断当前是否已经可以再次刷新。
func Ready(now, retryAt time.Time) bool {
	return !now.Before(retryAt)
}
