package vagons

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func HoldPanic(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "<p><font color=\"red\">", r, "</font></p>", "<p><code>",
				string(debug.Stack()), "</code></p>")
		}
	}()
	next(w, r)
}
