package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type GeoService struct {
	Rdb *redis.Client
}

const KeyDrivers = "drivers"

func (g *GeoService) UpdateLocation(ctx context.Context, driverId string, lat, long float64) error {
	_, err := g.Rdb.GeoAdd(ctx, KeyDrivers, &redis.GeoLocation{
		Name:      driverId,
		Longitude: long,
		Latitude:  lat,
	}).Result()
	return err
}

func (g *GeoService) FindNearestDrive(ctx context.Context, riderLat, riderLong float64) ([]redis.GeoLocation, error) {
	// Use GeoRadius (GeoRadiusOpts) with WithDist to get distance information.
	// Note: go-redis v9 provides GeoRadius with GeoRadiusQuery; use Radius to search and return distances.
	results, err := g.Rdb.GeoRadius(ctx, KeyDrivers, riderLong, riderLat, &redis.GeoRadiusQuery{
		Radius:   5,
		Unit:     "km",
		WithDist: true,
		Count:    1,
		Sort:     "ASC",
	}).Result()
	// GeoRadius returns a slice of GeoLocation which includes Dist (float64)
	return results, err
}
