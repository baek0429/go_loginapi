package login

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo?Username=baek&Password=test&Tag=Login", nil)
	if err != nil {
		t.Log(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)
	t.Logf("%d - %s", w.Code, w.Body.String())
}

func TestResponse(t *testing.T) {
}

// Hash test
func _TestGetStore(t *testing.T) {
	var a [10]string
	for i := 0; i < 10; i++ {
		a[i] = GetStoreUserQuery("helloworld", "hello")[2]
	}
	for j := 0; j < 10; j++ {
		t.Log(a[j])
		t.Log(ConfirmPassword(a[j], "hello"))
	}
}
