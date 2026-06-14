package handler

import (
	"net/http"
	"testing"

	"clawbench/internal/service"
)

func TestServeKeyConfig_GetEmpty(t *testing.T) {
	_, teardown := setupTestEnv(t)
	defer teardown()

	req := newRequest(t, http.MethodGet, "/api/terminal/key-config?type=key", nil)
	w := callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusOK)

	var items []service.KeyConfigItem
	decodeRespJSON(t, w.Body, &items)
	if len(items) != 0 {
		t.Fatalf("expected empty, got %d items", len(items))
	}
}

func TestServeKeyConfig_PutAndGet(t *testing.T) {
	_, teardown := setupTestEnv(t)
	defer teardown()

	// Put key config
	putBody := map[string]any{
		"type":  "key",
		"items": []string{"esc", "tab", "ctrl"},
	}
	req := newRequest(t, http.MethodPut, "/api/terminal/key-config", putBody)
	w := callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusOK)

	// Get key config
	req = newRequest(t, http.MethodGet, "/api/terminal/key-config?type=key", nil)
	w = callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusOK)

	var items []service.KeyConfigItem
	decodeRespJSON(t, w.Body, &items)
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].KeyID != "esc" || items[1].KeyID != "tab" || items[2].KeyID != "ctrl" {
		t.Fatalf("unexpected order: %+v", items)
	}
}

func TestServeKeyConfig_InvalidType(t *testing.T) {
	_, teardown := setupTestEnv(t)
	defer teardown()

	req := newRequest(t, http.MethodGet, "/api/terminal/key-config?type=invalid", nil)
	w := callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestServeKeyConfig_Replace(t *testing.T) {
	_, teardown := setupTestEnv(t)
	defer teardown()

	// Initial config
	putBody := map[string]any{
		"type":  "symbol",
		"items": []string{".", "/", "-"},
	}
	req := newRequest(t, http.MethodPut, "/api/terminal/key-config", putBody)
	w := callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusOK)

	// Replace with different config
	putBody = map[string]any{
		"type":  "symbol",
		"items": []string{"$", "&"},
	}
	req = newRequest(t, http.MethodPut, "/api/terminal/key-config", putBody)
	w = callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusOK)

	// Get should show only new config
	req = newRequest(t, http.MethodGet, "/api/terminal/key-config?type=symbol", nil)
	w = callHandler(ServeKeyConfig, req)
	assertStatus(t, w, http.StatusOK)

	var items []service.KeyConfigItem
	decodeRespJSON(t, w.Body, &items)
	if len(items) != 2 {
		t.Fatalf("expected 2 items after replace, got %d", len(items))
	}
	if items[0].KeyID != "$" || items[1].KeyID != "&" {
		t.Fatalf("unexpected items after replace: %+v", items)
	}
}
