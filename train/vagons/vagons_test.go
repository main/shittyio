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
		panic("manual panic")
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
	if string(body[0:21]) != "<p><font color=\"red\">" ||
		string(body[len(body)-10:]) != "/code></p>" {
		t.Log("body incorrect")
		t.FailNow()
	}
}
