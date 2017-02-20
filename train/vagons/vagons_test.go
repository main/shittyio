package vagons

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/main/shittyio/train"
)

func TestHoldPanic(t *testing.T) {

	trn := train.New(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("reached last handler")
		panic(1)
	}))
	trn.AddVagon(HoldPanic)
	req := httptest.NewRequest("GET", "/something", nil)
	w := httptest.NewRecorder()
	trn.Handler().ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Log("status code incorrect", resp.StatusCode)
		t.FailNow()
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "Panic" {
		t.Log("body incorrect")
		t.FailNow()
	}
}
