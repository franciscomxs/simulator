package lambdahandler

import (
	"math"
	"testing"
)

func TestJsonResponse_Success(t *testing.T) {
	resp := jsonResponse(200, map[string]string{"status": "ok"})
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected application/json, got %q", resp.Headers["Content-Type"])
	}
	if resp.Body != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", resp.Body)
	}
}

func TestJsonResponse_MarshalError(t *testing.T) {
	// json.Marshal fails on NaN — exercises the error branch
	type nanPayload struct{ V float64 }
	resp := jsonResponse(200, nanPayload{V: math.NaN()})
	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
	if resp.Body != `{"error":"internal server error"}` {
		t.Errorf("unexpected body: %s", resp.Body)
	}
}
