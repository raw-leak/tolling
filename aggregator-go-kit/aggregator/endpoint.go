package aggregator

import (
	"context"
	"errors"
	"log"
	"net/http"
	"tolling/types"

	"github.com/go-kit/kit/endpoint"
)

// Collects a set of endpoints
type Set struct {
	AggregateEndpoint  endpoint.Endpoint
	GetInvoiceEndpoint endpoint.Endpoint
}

func New(svc Service) {

	endpointSet := Set{
		AggregateEndpoint:  makeAggregateEndpoint(svc),
		GetInvoiceEndpoint: makeGetInvoiceEndpoint(svc),
	}

	handler := MakeHTTPHandler(endpointSet)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("HTTP Server Error: ", err)
	}
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

func (ar AggregateRequest) Validate() error {
	if ar.OBUID == 0 {
		return errors.New("OBUID is missing")
	}
	if ar.Value == 0.0 {
		return errors.New("Value is missing")
	}
	if ar.Unix == 0.0 {
		return errors.New("Unix is missing")
	}
	return nil
}

type AggregateResponse struct {
	Err error `json:"-"`
}

type GetInvoiceRequest struct {
	OBUID int `json:"obuID"`
}

func (gir GetInvoiceRequest) Validate() error {
	if gir.OBUID == 0 {
		return errors.New("OBUID is missing")
	}
	return nil
}

type GetInvoiceResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"-"`
}

func (s Set) Aggregate(ctx context.Context, d types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: d.OBUID,
		Value: d.Value,
		Unix:  d.Unix,
	})

	return err
}

func (s Set) GetInvoice(ctx context.Context, OBUDID int) (types.Invoice, error) {
	res, err := s.GetInvoiceEndpoint(ctx, GetInvoiceRequest{
		OBUID: OBUDID,
	})

	if err != nil {
		return nil, err
	}

	result := res.(GetInvoiceResponse)

	return types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func makeAggregateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		aggReq := request.(AggregateRequest)

		if err := aggReq.Validate(); err != nil {
			return AggregateResponse{Err: ValidationError{err.Error()}}, nil
		}

		err := s.Aggregate(ctx, types.Distance{
			OBUID: aggReq.OBUID,
			Value: aggReq.Value,
			Unix:  aggReq.Unix,
		})

		if err != nil {
			return nil, err
		}

		return AggregateResponse{Err: err}, nil
	}
}

func makeGetInvoiceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetInvoiceRequest)
		res, err := s.GetInvoice(ctx, req.OBUID)

		if err != nil {
			return nil, err
		}

		return GetInvoiceResponse{
			Err:           err,
			OBUID:         res.OBUID,
			TotalDistance: res.TotalDistance,
			TotalAmount:   res.TotalAmount,
		}, nil
	}
}
