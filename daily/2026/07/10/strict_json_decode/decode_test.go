package strictjsondecode

import "testing"

func TestDecodeRequestSuccess(t *testing.T) {
	req, err := DecodeRequest([]byte(`{"name":"worker-a","limit":8}`))
	if err != nil {
		t.Fatalf("DecodeRequest() error = %v", err)
	}
	if req.Name != "worker-a" || req.Limit != 8 {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestDecodeRequestRejectsUnknownField(t *testing.T) {
	_, err := DecodeRequest([]byte(`{"name":"worker-a","limit":8,"debug":true}`))
	if err == nil || err.Error() != `unknown field: json: unknown field "debug"` {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func TestDecodeRequestRejectsTrailingPayload(t *testing.T) {
	_, err := DecodeRequest([]byte(`{"name":"worker-a","limit":8} {"name":"worker-b","limit":4}`))
	if err == nil || err.Error() != "body must contain a single JSON object" {
		t.Fatalf("expected trailing payload error, got %v", err)
	}
}

func TestDecodeRequestRejectsInvalidBusinessFields(t *testing.T) {
	_, err := DecodeRequest([]byte(`{"name":" ","limit":0}`))
	if err == nil || err.Error() != "name is required" {
		t.Fatalf("expected business validation error, got %v", err)
	}
}
