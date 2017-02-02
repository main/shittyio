package train

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestTrainBasic(t *testing.T) {
	events := []string{}
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			t.Log("cliet reached /")
			w.WriteHeader(404)
			fmt.Fprintln(w, "You're lost, go home")
			events = append(events, "/")
		})
		mux.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
			t.Log("cliet reached /page1")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
			events = append(events, "/page1")
		})

		trn := New(mux)

		{
			vagon1 := func(w http.ResponseWriter, r *http.Request) {
				t.Log("vagon 1 launched")
				events = append(events, "vagon 1")
			}
			trn.AddVagon(vagon1)
		}

		{
			vagon2 := func(w http.ResponseWriter, r *http.Request) {
				t.Log("vagon 2 launched")
				events = append(events, "vagon 2")
			}
			trn.AddVagon(vagon2)
		}

		if err := http.ListenAndServe(":9999", trn); err != nil {
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
	body := new(bytes.Buffer)
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
	if !reflect.DeepEqual(events, []string{"vagon 1", "vagon 2", "/page1"}) {
		t.Log("events isn't the same like etalon", events)
		t.FailNow()
	}
}
