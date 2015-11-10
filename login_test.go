package login

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo?Username=baek&Password=test", nil)
	if err != nil {
		t.Log(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)
	t.Logf("%d - %s", w.Code, w.Body.String())
}
