package log_field_redactor

import "strings"

type Redactor struct{ secrets map[string]struct{} }

func New(secretFields ...string) Redactor {
	secrets := make(map[string]struct{}, len(secretFields))
	for _, field := range secretFields {
		secrets[strings.ToLower(field)] = struct{}{}
	}
	return Redactor{secrets: secrets}
}

func (r Redactor) Redact(fields map[string]any) map[string]any {
	result := make(map[string]any, len(fields))
	for key, value := range fields {
		lower := strings.ToLower(key)
		if _, secret := r.secrets[lower]; secret {
			continue
		}
		if lower == "email" {
			if email, ok := value.(string); ok {
				if at := strings.LastIndex(email, "@"); at > 0 && at < len(email)-1 {
					value = "***" + email[at:]
				}
			}
		}
		result[key] = value
	}
	return result
}
