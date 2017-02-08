package train

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type Events []string

func (events *Events) Handler(event string, status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*events = append(*events, event)
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}
}

func (events *Events) Vagon(event string) VagonFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		*events = append(*events, event)
		next(w, r)
	}
}

func (events Events) CompareWithEthalon(ethalon ...string) bool {
	return reflect.DeepEqual([]string(events), ethalon)
}

func TestTrainBasic(t *testing.T) {
	events := Events{}
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/page1", events.Handler("/page1", http.StatusOK, "OK"))
		trn := New(mux)
		{
			vagon1 := events.Vagon("vagon 1")
			trn.AddVagon(vagon1)
		}

		{
			vagon2 := events.Vagon("vagon 2")
			trn.AddVagon(vagon2)
		}

		if err := http.ListenAndServe(":9999", trn.Handler()); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}()
	response, err := http.Get("http://localhost:9999/page1")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusOK {
		t.Log(response.Status)
		t.FailNow()
	}
	defer response.Body.Close()
	//	bufer := make([]byte, 100)
	//	_, err = response.Body.Read(bufer)

	//	body, err := ioutil.ReadAll(response.Body)
	body := bytes.Buffer{}
	_, err = body.ReadFrom(response.Body)
	if err != nil && err != io.EOF {
		t.Log(err)
		t.FailNow()
	}
	bodyVampire := body.String()
	if bodyVampire != "OK" {
		t.Log(fmt.Sprintf("%#v", bodyVampire))
		t.FailNow()
	}
	if !events.CompareWithEthalon("vagon 1", "vagon 2", "/page1") {
		t.Log("events isn't the same like etalon", events)
		t.FailNow()
	}
}

func TestTrainReject(t *testing.T) {
	events := Events{}
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/page1", events.Handler("/page1", http.StatusOK, "OK"))
		trn := New(mux)
		{
			vagon1 := events.Vagon("vagon 1")
			trn.AddVagon(vagon1)
		}
		{
			vagon2 := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
				events = append(events, "vagon 2")
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, "Forbidden")
			}
			trn.AddVagon(vagon2)
		}
		{
			vagon3 := events.Vagon("vagon 3")
			trn.AddVagon(vagon3)
		}

		if err := http.ListenAndServe(":9998", trn.Handler()); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}()
	response, err := http.Get("http://localhost:9998/page1")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if response.StatusCode != http.StatusForbidden {
		t.Log(response.Status)
		t.FailNow()
	}
	defer response.Body.Close()

	body := bytes.Buffer{}
	_, err = body.ReadFrom(response.Body)
	if err != nil && err != io.EOF {
		t.Log(err)
		t.FailNow()
	}
	bodyVampire := body.String()
	if bodyVampire != "Forbidden" {
		t.Log(fmt.Sprintf("%#v", bodyVampire))
		t.FailNow()
	}
	if !events.CompareWithEthalon("vagon 1", "vagon 2") {
		t.Log("events isn't the same like etalon", events)
		t.FailNow()
	}
}
