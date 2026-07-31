package problemdetailsanitizer

import "strings"

// Problem 表示 RFC 9457 风格的问题详情。
type Problem struct {
	Type     string
	Title    string
	Status   int
	Detail   string
	Instance string
}

// Sanitize 仅允许预先登记的公开详情离开服务边界。
func Sanitize(problem Problem, allowedDetails []string) Problem {
	result := Problem{
		Type:   problem.Type,
		Title:  problem.Title,
		Status: problem.Status,
		Detail: "request could not be completed",
	}
	if result.Status < 400 || result.Status > 599 {
		result.Status = 500
	}
	for _, allowed := range allowedDetails {
		if allowed != "" && strings.Contains(problem.Detail, allowed) {
			result.Detail = allowed
			break
		}
	}
	return result
}
