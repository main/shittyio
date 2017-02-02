package train

import "net/http"

type VagonFunc func(w http.ResponseWriter, r *http.Request)

var _ http.Handler = new(Train)

type Train struct {
	handler http.Handler
	vagons  []VagonFunc
}

func New(handler http.Handler) *Train {
	return &Train{handler: handler}
}

func (train *Train) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, vagon := range train.vagons {
		vagon(w, r)
	}
	train.handler.ServeHTTP(w, r)
}

func (train *Train) AddVagon(vagon VagonFunc) {
	train.vagons = append(train.vagons, vagon)
}
