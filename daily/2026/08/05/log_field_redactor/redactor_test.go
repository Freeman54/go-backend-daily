package log_field_redactor

import "testing"

func TestRedact_RemovesSecretsAndMasksEmail(t *testing.T) {
	r := New("password", "authorization")
	input := map[string]any{"Password": "secret", "email": "alice@example.com", "request_id": "r1"}
	got := r.Redact(input)
	if _, ok := got["Password"]; ok {
		t.Fatal("密码未删除")
	}
	if got["email"] != "***@example.com" {
		t.Fatalf("email = %v", got["email"])
	}
	if got["request_id"] != "r1" {
		t.Fatalf("request_id = %v", got["request_id"])
	}
	if input["Password"] != "secret" {
		t.Fatal("输入被修改")
	}
}

func TestRedact_LeavesInvalidEmailUnchanged(t *testing.T) {
	got := New().Redact(map[string]any{"email": "unknown"})
	if got["email"] != "unknown" {
		t.Fatalf("email = %v", got["email"])
	}
}
