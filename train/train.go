package train

import "net/http"

type VagonFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type Train struct {
	handler http.Handler
	vagons  []VagonFunc
}

func New(handler http.Handler) *Train {
	return &Train{handler: handler}
}

func (train *Train) Handler() http.Handler {
	compose := func(v VagonFunc, next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			v(w, r, next)
		}
	}
	lastVagon := train.handler.ServeHTTP
	for i := len(train.vagons) - 1; i >= 0; i-- {
		lastVagon = compose(train.vagons[i], lastVagon)
	}
	return http.HandlerFunc(lastVagon)
}

func (train *Train) AddVagon(vagon VagonFunc) {
	train.vagons = append(train.vagons, vagon)
}
