package attributeallowlist

// Filter 仅保留允许跨服务边界传播的消息属性。
// 返回新 map，避免调用方后续修改原始消息时影响过滤结果。
func Filter(attributes map[string]string, allowed []string) map[string]string {
	allow := make(map[string]struct{}, len(allowed))
	for _, key := range allowed {
		allow[key] = struct{}{}
	}

	result := make(map[string]string)
	for key, value := range attributes {
		if _, ok := allow[key]; ok {
			result[key] = value
		}
	}
	return result
}
