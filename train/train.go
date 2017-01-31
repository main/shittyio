package train

import "net/http"

type VagonFunc func(w http.ResponseWriter, r *http.Request)

var _ http.Handler = new(Train)

type Train struct {
	handler http.Handler
	vagons  []VagonFunc
}

func (train *Train) New(handler http.Handler) *Train {
	train := &Train{}
	return train
}

func (train *Train) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (train *Train) AddVagon(vagon VagonFunc) {
	train.vagons = append(train.vagons, vagon)
}
