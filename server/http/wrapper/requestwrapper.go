package wrapper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HandlerError struct {
	Err        error
	StatusCode int
}

const LoggerCtxKey = "loggerCtxKey"
const URLParamsCtxKey = "urlParamsCtxKey"

type HTTPWrapper struct {
	logger          logrus.FieldLogger
	PossibleMethods []string
}

func NewHTTPWrapper(
	logger logrus.FieldLogger,
) *HTTPWrapper {
	httpWrapper := &HTTPWrapper{
		logger:          logger,
		PossibleMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
	}
	return httpWrapper
}

func (w *HTTPWrapper) Register(router *mux.Router, path string, handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	for i := 0; i < handlerType.NumMethod(); i++ {
		method := handlerType.Method(i)
		verb, ok := w.getHTTPVerb(method)

		if fitsHandlerTemplate(method) && ok {
			w.logger.Infof("Registered %s %s %s.%s", verb, path, handlerType, method.Name)
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
			"handler":    handler.Type(),
			"methodName": method.Name,
		})

		payload, varsValue, err := unmarshallArgsAndVars(method, r)
		if err != nil {
			logger.WithError(err).Error("Failed to unmarshall request")
			writter.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, LoggerCtxKey, logger)
		ctx = context.WithValue(ctx, URLParamsCtxKey, r.URL.Query())
		args := payload.Addr().Interface()
		vars := varsValue.Addr().Interface()

		responseSlice := method.Func.Call([]reflect.Value{
			handler, reflect.ValueOf(ctx), reflect.ValueOf(args), reflect.ValueOf(vars),
		})

		handlerResult, handlerErr := responseSlice[0], responseSlice[1]
		if !handlerErr.IsNil() {
			writeErr(writter, logger, handlerErr)
			return
		}

		bytes, err := json.Marshal(handlerResult.Interface())
		if err != nil {
			logger.WithError(err).Error("Could not unmarshall handler response")
			writter.WriteHeader(http.StatusInternalServerError)
			return
		}

		writter.WriteHeader(http.StatusOK)
		writter.Write(bytes)
		logger.Debug("Request processed")
	}
}

func unmarshallArgsAndVars(method reflect.Method, r *http.Request) (reflect.Value, reflect.Value, error) {
	defer r.Body.Close()

	payloadType := method.Type.In(2)
	payload := instantiate(payloadType)

	varsType := method.Type.In(3)
	varsValue := instantiate(varsType)

	err := json.NewDecoder(r.Body).Decode(payload.Addr().Interface())
	switch {
	case err == io.EOF:
		// empty body, do nothing
	case err != nil:
		return payload, varsValue, err
	}

	muxVars := mux.Vars(r)
	if len(muxVars) > 0 {
		marshalledVars, err := json.Marshal(mux.Vars(r))
		json.Unmarshal(marshalledVars, varsValue.Addr().Interface())
		if err != nil {
			return payload, varsValue, err
		}
	}
	return payload, varsValue, nil
}

func writeErr(writter http.ResponseWriter, logger logrus.FieldLogger, handlerErrVal reflect.Value) {
	handlerEr := handlerErrVal.Elem().Interface().(HandlerError)

	logger.WithError(handlerEr.Err).Error("Handler returned error")
	writter.WriteHeader(handlerEr.StatusCode)
}

func instantiate(instanceType reflect.Type) reflect.Value {
	var instance reflect.Value
	if instanceType.Kind() == reflect.Ptr {
		instance = reflect.New(instanceType.Elem()).Elem()
	} else {
		instance = reflect.New(instanceType).Elem()
	}
	return instance
}

type HandlerTemplate struct{}

func (handlerTemplate *HandlerTemplate) RequestTemplate(ctx context.Context, pointer *struct{}, routeVars *struct{}) (*struct{}, *HandlerError) {
	return nil, nil
}

func fitsHandlerTemplate(method reflect.Method) bool {
	template := reflect.TypeOf(&HandlerTemplate{}).Method(0)

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
