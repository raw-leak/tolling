package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"tolling/common"
	"tolling/types"
)

type Aggregator interface {
	AggregateInvoice(ctx context.Context, d types.Distance) error
}

type Consumer[T any] interface {
	Consume(context.Context) chan EventPayload[T]
}

type Point struct {
	Long float64
	Lat  float64
}

type CalculatedDistance struct {
	mu            sync.Mutex
	lastDistance  float64
	pointsHistory []Point
}

type DistanceCalculator struct {
	con        Consumer[types.OBUData]
	distances  map[int]*CalculatedDistance
	aggregator Aggregator
	logger     common.Logger
}

func (dc *DistanceCalculator) Run(ctx context.Context) {
	consCh := dc.con.Consume(ctx)
	dc.logger.New().Info("started distance calculation process")

	// TODO: improve by spanning multiple go routines
	for p := range consCh {
		l := dc.logger.New()
		if traceID, ok := p.ctx.Value(types.KeyTraceID).(string); ok {
			fmt.Println("before sending to AggregateInvoice HTTP traceID", traceID, ok)
			l.WithTraceID(traceID)
		}

		dist := dc.calculateDistanceAndSavePoint(p.data)
		err := dc.aggregator.AggregateInvoice(p.ctx, types.Distance{OBUID: p.data.OBUID, Value: dist, Unix: time.Now().Unix()})

		if err != nil {
			l.WithError(err).Error("failed to send to aggregate invoice")
		}

	}

	dc.logger.New().Info("finished distance calculation process")
}

func NewDistanceCalculator(con Consumer[types.OBUData], aggregator Aggregator, logger common.Logger) *DistanceCalculator {
	return &DistanceCalculator{con, make(map[int]*CalculatedDistance), aggregator, logger}
}

func (ds *DistanceCalculator) calculateDistance(x1, x2, y1, y2 float64) float64 {
	return 0.0
}

func (ds *DistanceCalculator) calculateDistanceAndSavePoint(data types.OBUData) float64 {
	newPoint := Point{Long: data.Long, Lat: data.Lat}
	distance := 0.0

	calculatedDistance, ok := ds.distances[data.OBUID]
	if !ok {
		calculatedDistance = &CalculatedDistance{}
		ds.distances[data.OBUID] = calculatedDistance
	}

	calculatedDistance.mu.Lock()
	defer calculatedDistance.mu.Unlock()

	if len(calculatedDistance.pointsHistory) > 0 {
		prevPoint := calculatedDistance.pointsHistory[len(calculatedDistance.pointsHistory)-1]
		distance = ds.calculateDistance(newPoint.Lat, newPoint.Long, prevPoint.Lat, prevPoint.Long)
	}
	calculatedDistance.lastDistance = distance
	calculatedDistance.pointsHistory = append(calculatedDistance.pointsHistory, newPoint)

	return distance
}
