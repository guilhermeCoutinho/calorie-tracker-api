package http

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTPWrapper struct {
	logger          logrus.FieldLogger
	PossibleMethods []string
}

func NewHTTPWrapper(
	logger logrus.FieldLogger,
) *HTTPWrapper {
	httpWrapper := &HTTPWrapper{
		logger:          logger,
		PossibleMethods: []string{"GET", "POST"},
	}
	return httpWrapper
}

func (w *HTTPWrapper) Register(router *mux.Router, path string, handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	for i := 0; i < handlerType.NumMethod(); i++ {
		method := handlerType.Method(i)
		verb, ok := w.getHTTPVerb(method)
		if fitsHandlerTemplate(method) && ok {
			w.logger.Infof("Registered %s into %s/%s", method.Name, path, verb)
			router.HandleFunc(path, w.wrapHTTPRequest(reflect.ValueOf(handler), method)).Methods(verb)
		}
	}
}

func (w *HTTPWrapper) getHTTPVerb(method reflect.Method) (string, bool) {
	for _, verb := range w.PossibleMethods {
		if strings.HasPrefix(strings.ToLower(method.Name), strings.ToLower(verb)) {
			return verb, true
		}
	}
	return "", false
}

func (w *HTTPWrapper) wrapHTTPRequest(handler reflect.Value, method reflect.Method) func(w http.ResponseWriter, r *http.Request) {
	return func(writter http.ResponseWriter, r *http.Request) {
		logger := w.logger.WithFields(logrus.Fields{
			"methodName": method.Name,
		})

		requestType := method.Type.In(2)
		var instance reflect.Value
		if requestType.Kind() == reflect.Ptr {
			instance = reflect.New(requestType.Elem()).Elem()
		} else {
			instance = reflect.New(requestType).Elem()
		}

		err := json.NewDecoder(r.Body).Decode(instance.Addr().Interface())
		if err != nil {
			logger.WithError(err).Error("Failed to parse payload")
			writter.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		args := instance.Addr().Interface()
		responseSlice := method.Func.Call([]reflect.Value{
			handler, reflect.ValueOf(ctx), reflect.ValueOf(args),
		})

		handlerResult, handlerErr := responseSlice[0].Interface(), responseSlice[1].Interface().(error)

		if handlerErr != nil {
			logger.WithError(handlerErr).Error("Handler returned error")
			writter.WriteHeader(http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(handlerResult)
		if err != nil {
			logger.WithError(err).Error("Could not unmarshall handler response")
			writter.WriteHeader(http.StatusInternalServerError)
			return
		}

		writter.WriteHeader(http.StatusOK)
		writter.Write(bytes)
		logger.Info("finished processing request")
	}
}

type HandlerTemplate struct {
}

func (handlerTemplate *HandlerTemplate) PitayaRequestTemplate(ctx context.Context, pointer *struct{}) (*struct{}, error) {
	return nil, nil
}

func fitsHandlerTemplate(method reflect.Method) bool {
	template, ok := reflect.TypeOf(&HandlerTemplate{}).MethodByName("PitayaRequestTemplate")
	if !ok {
		return false
	}

	if method.Type.NumIn() != template.Type.NumIn() || method.Type.NumOut() != template.Type.NumOut() {
		return false
	}
	for i := 0; i < method.Type.NumIn(); i++ {
		if template.Type.In(i).Kind() != method.Type.In(i).Kind() {
			return false
		}
	}

	for i := 0; i < method.Type.NumOut(); i++ {
		if template.Type.Out(i).Kind() != method.Type.Out(i).Kind() {
			return false
		}
	}
	return true
}
