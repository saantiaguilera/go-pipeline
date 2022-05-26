package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	// Repositories is a grouping of all repositories in this sample
	// This is simply for showcase purposes
	Repositories struct {
		Driver        *DriverRepository
		Event         *EventRepository
		Location      *LocationRepository
		Notifications *NotificationRepository
		Tracking      *TrackingRepository
	}

	// Steps is a grouping of all steps in this sample
	// This is simply for showcase purposes
	Steps struct {
		GetDriverFromEvent        pipeline.Step[EventID, DriverID]
		GetDriverByID             pipeline.Step[DriverID, Driver]
		GetCurrentCoordinatesByID pipeline.Step[DriverID, Coordinate]
		GetLocation               pipeline.Step[Coordinate, Location]
		LocationToGeoDriver       pipeline.Step[Location, GeoDriver]
		DriverToGeoDriver         pipeline.Step[Driver, GeoDriver]
		NotifyClose               pipeline.Step[GeoDriver, GeoDriver]
		SaveInDestination         pipeline.Step[GeoDriver, GeoDriver]
		MarkProcessed             pipeline.Step[EventID, EventID]
	}

	// Stmts is a grouping of all statements in this sample
	// This is simply for showcase purposes
	Stmts struct {
		IsDriverClose pipeline.Statement[GeoDriver]
	}
)

var (
	render = flag.Bool("pipeline.render", false, "render pipeline")
)

// NewRepositories creates all the repositories for this sample.
// Since they are all dummies, we dont have to create anything
func NewRepositories() Repositories {
	return Repositories{}
}

// NewSteps creates all the steps for this sample
func NewSteps(r Repositories) Steps {
	return Steps{
		GetDriverFromEvent:        pipeline.NewUnitStep("unmarshall_event", r.Event.GetTrigger),
		GetDriverByID:             pipeline.NewUnitStep("get_driver", r.Driver.GetDriverByID),
		GetCurrentCoordinatesByID: pipeline.NewUnitStep("get_current_coordinates", r.Driver.GetRealTimeCoordinatesByID),
		GetLocation:               pipeline.NewUnitStep("get_location", r.Location.GetFullLocationByLatLng),
		LocationToGeoDriver: pipeline.NewUnitStep(
			"location_to_geodriver",
			func(ctx context.Context, l Location) (GeoDriver, error) {
				return GeoDriver{Location: l}, nil
			},
		),
		DriverToGeoDriver: pipeline.NewUnitStep(
			"driver_to_geodriver",
			func(ctx context.Context, d Driver) (GeoDriver, error) {
				return GeoDriver{Driver: d}, nil
			},
		),
		NotifyClose: pipeline.NewUnitStep(
			"notify_driver_close",
			func(ctx context.Context, gd GeoDriver) (GeoDriver, error) {
				return gd, r.Notifications.NotifyCloseToDestination(ctx, gd.Driver)
			},
		),
		SaveInDestination: pipeline.NewUnitStep(
			"save_driver_in_destination",
			func(ctx context.Context, gd GeoDriver) (GeoDriver, error) {
				return gd, r.Driver.SaveDriverInDestination(ctx, gd)
			},
		),
		MarkProcessed: pipeline.NewUnitStep(
			"mark_event_as_processed",
			func(ctx context.Context, ei EventID) (EventID, error) {
				return ei, r.Event.MarkProcessed(ctx, ei)
			},
		),
	}
}

// NewStmts creates all the statements for this sample
func NewStmts(r Repositories) Stmts {
	return Stmts{
		IsDriverClose: pipeline.NewStatement(
			"is_driver_close_to_destination",
			r.Tracking.IsDriverClose,
		),
	}
}

// NewGraph creates the pipeline for this sample
func NewGraph() pipeline.Step[EventID, GeoDriver] {
	r := NewRepositories()
	s := NewSteps(r)
	ss := NewStmts(r)

	// Pipeline creation. The pipeline looks like this:
	// event -> driver info         -> optional (if close -> notify      ) -> mark processed
	//       -> coords -> location              (         -> save close  )
	return NewEventStep[GeoDriver](
		pipeline.NewSequentialStep[EventID, GeoDriver](
			pipeline.NewSequentialStep(
				s.GetDriverFromEvent,
				newConcurrentGeoDriver(s),
			),
			newProcessDriverCloseToDestination(s, ss),
		),
		s.MarkProcessed,
	)
}

// RunPipeline runs the provided pipeline.
// Output: (one of many)
//
// getting driver id by event id 1234
// getting driver real time coordinates by id 1234
// getting driver data by id 1234
// getting full location by coordinates {Lat:32 Lng:1}
// checking if driver is close to destination {Driver:{ID:1234} Location:{Coordinate:{Lat:32 Lng:1} City:example State:example Country:example}}
// driver is close to destination!
// saving driver in destination {Driver:{ID:1234} Location:{Coordinate:{Lat:32 Lng:1} City:example State:example Country:example}}
// sending notification to driver {ID:1234}
// marking event id 1234 as processed
// {Driver:{ID:1234} Location:{Coordinate:{Lat:32 Lng:1} City:example State:example Country:example}}
func RunPipeline() {
	// Create a stateless graph, so it can be evaluated as many times as we like with any input/context we want to.
	pipe := NewGraph()

	// Initial input data
	ctx := context.Background()
	in := EventID(1234)

	// Run and assert.
	res, err := pipe.Run(ctx, in)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}

func RunGraphRendering() {
	if *render {
		diagram := pipeline.NewUMLGraph()
		renderer := pipeline.NewUMLRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("template.svg")

		NewGraph().Draw(diagram)

		if err := renderer.Render(diagram, file); err != nil {
			panic(err)
		}
	}
}

// main runs the sample. Note how in this sample none of the code is coupled to the pipeline package
// (besides the main one, which should be the cmd stuff that injects everything and of course is
// coupled to implementations)
//
// Since the API has an extremely flexible contract, there shouldn't be any need to couple repositories
// or structures to the API and they can simply behave normally as they would expect to.
func main() {
	RunGraphRendering()
	RunPipeline()
}

func newGeoDriverLocationStep(s Steps) pipeline.Step[DriverID, GeoDriver] {
	getFullLocation := pipeline.NewSequentialStep(s.GetCurrentCoordinatesByID, s.GetLocation)
	return pipeline.NewSequentialStep[DriverID, Location](getFullLocation, s.LocationToGeoDriver)
}

func newGeoDriverDriverStep(s Steps) pipeline.Step[DriverID, GeoDriver] {
	return pipeline.NewSequentialStep(s.GetDriverByID, s.DriverToGeoDriver)
}

func newConcurrentGeoDriver(s Steps) pipeline.Step[DriverID, GeoDriver] {
	return pipeline.NewConcurrentStep(
		[]pipeline.Step[DriverID, GeoDriver]{
			newGeoDriverDriverStep(s),
			newGeoDriverLocationStep(s),
		},
		func(ctx context.Context, gd1, gd2 GeoDriver) (GeoDriver, error) {
			if gd2.Location != (Location{}) {
				gd1.Location = gd2.Location
			}
			if gd2.Driver != (Driver{}) {
				gd1.Driver = gd2.Driver
			}
			return gd1, nil
		},
	)
}

func newConcurrentCloseToDestination(s Steps) pipeline.Step[GeoDriver, GeoDriver] {
	return pipeline.NewConcurrentStep(
		[]pipeline.Step[GeoDriver, GeoDriver]{
			s.NotifyClose,
			s.SaveInDestination,
		},
		func(ctx context.Context, gd1, gd2 GeoDriver) (GeoDriver, error) {
			return gd1, nil // no need to reduce as they don't produce output.
		},
	)
}

func newProcessDriverCloseToDestination(s Steps, ss Stmts) pipeline.Step[GeoDriver, GeoDriver] {
	return pipeline.NewOptionalStep(ss.IsDriverClose, newConcurrentCloseToDestination(s))
}
