package pipe

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	hub "github.com/thingful/device-hub"
)

func NewHTTPListener(config map[string]interface{}) (*httpListener, error) {

	binding, found := config["HTTPBindingAddress"]

	if !found {
		return nil, errors.New("unable to find binding in configuration")
	}

	router := DefaultRouter()

	startDefaultHTTPListener(router, binding.(string))

	return &httpListener{
		router: router,
	}, nil
}

type httpListener struct {
	router *router
}

func (h *httpListener) NewChannel(uri string) (Channel, error) {

	errors := make(chan error)
	out := make(chan hub.Input)

	channel := defaultChannel{out: out, errors: errors}

	h.router.register(uri, channel)
	return channel, nil
}

func (h *httpListener) Close() error {
	return nil
}

func startDefaultHTTPListener(router *router, binding string) {

	http.HandleFunc("/", rootHandler(router))

	// TODO : shutdown nicely
	go func() {
		log.Fatal(http.ListenAndServe(binding, nil))

	}()
}

func rootHandler(router *router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		path := r.URL.Path

		ok, channel := router.Match(path)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			channel.Errors() <- err
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		input := hub.Input{
			Payload: body,
		}

		channel.Out() <- input
		w.WriteHeader(http.StatusAccepted)
	}
}
