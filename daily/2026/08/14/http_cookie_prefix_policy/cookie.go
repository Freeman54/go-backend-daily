package cookieprefix

import (
	"fmt"
	"net/http"
	"strings"
)

// Validate 检查 __Secure- 和 __Host- Cookie 的安全属性。
func Validate(cookie *http.Cookie) error {
	if cookie == nil {
		return fmt.Errorf("cookie 不能为空")
	}
	if strings.HasPrefix(cookie.Name, "__Secure-") && !cookie.Secure {
		return fmt.Errorf("__Secure- Cookie 必须设置 Secure")
	}
	if strings.HasPrefix(cookie.Name, "__Host-") {
		if !cookie.Secure {
			return fmt.Errorf("__Host- Cookie 必须设置 Secure")
		}
		if cookie.Path != "/" {
			return fmt.Errorf("__Host- Cookie 的 Path 必须为 / ")
		}
		if cookie.Domain != "" {
			return fmt.Errorf("__Host- Cookie 不能设置 Domain")
		}
	}
	return nil
}
