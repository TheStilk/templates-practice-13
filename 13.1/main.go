package main

import (
	"errors"
	"fmt"
	"time"
)

type RideState string

const (
	StateIdle           RideState = "Idle"
	StateCarSelected    RideState = "CarSelected"
	StateOrderConfirmed RideState = "OrderConfirmed"
	StateCarArrived     RideState = "CarArrived"
	StateInTrip         RideState = "InTrip"
	StateTripCompleted  RideState = "TripCompleted"
	StateTripCancelled  RideState = "TripCancelled"
)

type RideOrder struct {
	ID     string
	State  RideState
	CarID  string
	Driver string
	Rating int
}

type RideEvent string

const (
	EventSelectCar       RideEvent = "selectCar"
	EventConfirmOrder    RideEvent = "confirmOrder"
	EventCarArrived      RideEvent = "carArrived"
	EventStartTrip       RideEvent = "startTrip"
	EventEndTrip         RideEvent = "endTrip"
	EventCancelOrder     RideEvent = "cancelOrder"
	EventCarDelayed      RideEvent = "carDelayed"
	EventPaymentSuccess  RideEvent = "paymentSuccess"
	EventPaymentFailed   RideEvent = "paymentFailed"
	EventChangeCar       RideEvent = "changeCar"
	EventEmergencyCancel RideEvent = "emergencyCancel"
)

var transitions = map[RideState]map[RideEvent]RideState{
	StateIdle: {
		EventSelectCar:   StateCarSelected,
		EventCancelOrder: StateTripCancelled,
	},
	StateCarSelected: {
		EventConfirmOrder: StateOrderConfirmed,
		EventChangeCar:    StateCarSelected,
		EventCancelOrder:  StateTripCancelled,
	},
	StateOrderConfirmed: {
		EventCarArrived:  StateCarArrived,
		EventCancelOrder: StateTripCancelled,
		EventCarDelayed:  StateTripCancelled,
	},
	StateCarArrived: {
		EventStartTrip:   StateInTrip,
		EventCancelOrder: StateTripCancelled,
	},
	StateInTrip: {
		EventEndTrip:         StateTripCompleted,
		EventEmergencyCancel: StateTripCancelled,
	},
	StateTripCompleted: {
		EventPaymentSuccess: StateIdle,
		EventPaymentFailed:  StateTripCompleted,
	},
	StateTripCancelled: {},
}

func (r *RideOrder) CanTransition(event RideEvent) bool {
	_, ok := transitions[r.State][event]
	return ok
}

func (r *RideOrder) Transition(event RideEvent) error {
	if !r.CanTransition(event) {
		return fmt.Errorf("invalid transition: %s -> %s", r.State, event)
	}
	newState := transitions[r.State][event]
	fmt.Printf("Order %s: %s -> %s\n", r.ID, r.State, newState)
	r.State = newState

	switch event {
	case EventSelectCar:
		fmt.Println("Car selected.")
	case EventConfirmOrder:
		fmt.Println("Order confirmed. Car is on the way.")
	case EventCarArrived:
		fmt.Println("Car has arrived.")
	case EventStartTrip:
		fmt.Println("Trip started.")
	case EventEndTrip:
		fmt.Println("Trip completed. Payment pending.")
	case EventCancelOrder, EventCarDelayed, EventEmergencyCancel:
		fmt.Println("Order cancelled.")
	case EventPaymentSuccess:
		fmt.Println("Payment successful.")
	case EventPaymentFailed:
		fmt.Println("Payment failed. Please try again.")
	}

	return nil
}

func (r *RideOrder) SimulateDelay() {
	if r.State == StateOrderConfirmed {
		time.Sleep(2 * time.Second) // simulate waiting
		fmt.Println("Car is delayed...")
		r.Transition(EventCarDelayed)
	}
}

func (r *RideOrder) SubmitRating(rating int) error {
	if r.State != StateIdle {
		return errors.New("rating can only be submitted after the trip cycle is complete")
	}
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	r.Rating = rating
	fmt.Printf("Thank you for the rating: %d\n", rating)
	return nil
}

func main() {
	order := &RideOrder{
		ID:    "RIDE-001",
		State: StateIdle,
	}

	order.Transition(EventSelectCar)
	order.Transition(EventConfirmOrder)
	order.Transition(EventCarArrived)
	order.Transition(EventStartTrip)
	order.Transition(EventEndTrip)
	order.Transition(EventPaymentSuccess)

	order.SubmitRating(5)

	fmt.Println("\n--- Scenario with cancellation ---")
	order2 := &RideOrder{ID: "RIDE-002", State: StateIdle}
	order2.Transition(EventSelectCar)
	order2.Transition(EventCancelOrder)

	fmt.Println("\n--- Scenario with delay ---")
	order3 := &RideOrder{ID: "RIDE-003", State: StateIdle}
	order3.Transition(EventSelectCar)
	order3.Transition(EventConfirmOrder)
	go order3.SimulateDelay()
	time.Sleep(3 * time.Second)
}
