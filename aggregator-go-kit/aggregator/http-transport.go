package aggregator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// https://github.com/go-kit/examples/blob/master/addsvc/pkg/addtransport/http.go

func MakeHTTPHandler(s Set) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		// httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r.Methods("POST").Path("/aggregate").Handler(httptransport.NewServer(
		s.AggregateEndpoint,
		decodeAggregateRequest,
		encodeHTTPGenericResponse,
		options...,
	))

	r.Methods("GET").Path("/invoice").Handler(httptransport.NewServer(
		s.GetInvoiceEndpoint,
		decodeGetInvoiceRequest,
		encodeHTTPGenericResponse,
		options...,
	))

	return r
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAggregateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AggregateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetInvoiceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetInvoiceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(errToCode(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func errToCode(err error) int {
	// switch err {
	// case addservice.ErrTwoZeroes, addservice.ErrMaxSizeExceeded, addservice.ErrIntOverflow:
	// 	return http.StatusBadRequest
	// }

	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Message string
}

func (ve ValidationError) Error() string {
	return ve.Message
}
