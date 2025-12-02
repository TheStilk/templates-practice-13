package main

import (
	"fmt"
	"time"
)

type Role string

const (
	RoleGuest Role = "guest"
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID   int
	Name string
	Role Role
}

type Event struct {
	ID    int
	Title string
	Date  time.Time
	Venue string
}

type BookingStatus string

const (
	StatusActive    BookingStatus = "active"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID     int
	User   *User
	Event  *Event
	Status BookingStatus
}

type BookingSystem struct {
	events        []*Event
	users         []*User
	bookings      []*Booking
	nextEventID   int
	nextBookingID int
}

func NewBookingSystem() *BookingSystem {
	return &BookingSystem{
		events:        make([]*Event, 0),
		users:         make([]*User, 0),
		bookings:      make([]*Booking, 0),
		nextEventID:   1,
		nextBookingID: 1,
	}
}

func (s *BookingSystem) AddEvent(title string, date time.Time, venue string, admin *User) error {
	if admin.Role != RoleAdmin {
		return fmt.Errorf("only admin can add events")
	}
	event := &Event{
		ID:    s.nextEventID,
		Title: title,
		Date:  date,
		Venue: venue,
	}
	s.events = append(s.events, event)
	s.nextEventID++
	fmt.Printf("Event '%s' added (ID: %d)\n", title, event.ID)
	return nil
}

func (s *BookingSystem) UpdateEvent(eventID int, title string, date time.Time, venue string, admin *User) error {
	if admin.Role != RoleAdmin {
		return fmt.Errorf("only admin can edit events")
	}
	for _, e := range s.events {
		if e.ID == eventID {
			e.Title = title
			e.Date = date
			e.Venue = venue
			fmt.Printf("Event ID %d updated\n", eventID)
			return nil
		}
	}
	return fmt.Errorf("event not found")
}

func (s *BookingSystem) DeleteEvent(eventID int, admin *User) error {
	if admin.Role != RoleAdmin {
		return fmt.Errorf("only admin can delete events")
	}
	for i, e := range s.events {
		if e.ID == eventID {
			s.events = append(s.events[:i], s.events[i+1:]...)
			fmt.Printf("Event ID %d deleted\n", eventID)
			return nil
		}
	}
	return fmt.Errorf("event not found")
}

func (s *BookingSystem) ListEvents() {
	if len(s.events) == 0 {
		fmt.Println("No events available")
		return
	}
	fmt.Println("\nAvailable events:")
	for _, e := range s.events {
		fmt.Printf("ID: %d | %s | %s | %s\n",
			e.ID, e.Title, e.Date.Format("2006-01-02 15:04"), e.Venue)
	}
}

func (s *BookingSystem) BookEvent(userID, eventID int, user *User) error {
	if user.Role != RoleUser {
		return fmt.Errorf("only registered users can book")
	}
	var targetEvent *Event
	for _, e := range s.events {
		if e.ID == eventID {
			targetEvent = e
			break
		}
	}
	if targetEvent == nil {
		return fmt.Errorf("event not found")
	}
	booking := &Booking{
		ID:     s.nextBookingID,
		User:   user,
		Event:  targetEvent,
		Status: StatusActive,
	}
	s.bookings = append(s.bookings, booking)
	s.nextBookingID++
	fmt.Printf("Booking created: %s -> %s (ID: %d)\n", user.Name, targetEvent.Title, booking.ID)
	return nil
}

func (s *BookingSystem) CancelBooking(bookingID int, user *User) error {
	for _, b := range s.bookings {
		if b.ID == bookingID {
			if b.User.ID != user.ID && user.Role != RoleAdmin {
				return fmt.Errorf("you can only cancel your own bookings")
			}
			b.Status = StatusCancelled
			fmt.Printf("Booking ID %d cancelled\n", bookingID)
			return nil
		}
	}
	return fmt.Errorf("booking not found")
}

func (s *BookingSystem) ListAllBookings(admin *User) {
	if admin.Role != RoleAdmin {
		fmt.Println("Access denied")
		return
	}
	fmt.Println("\nAll bookings:")
	for _, b := range s.bookings {
		fmt.Printf("ID: %d | User: %s | Event: %s | Status: %s\n",
			b.ID, b.User.Name, b.Event.Title, b.Status)
	}
}

func main() {
	system := NewBookingSystem()

	guest := &User{ID: 1, Name: "Anna (guest)", Role: RoleGuest}
	user := &User{ID: 2, Name: "Ivan (user)", Role: RoleUser}
	admin := &User{ID: 3, Name: "Olga (admin)", Role: RoleAdmin}

	system.AddEvent("Jazz Concert", time.Now().Add(24*time.Hour), "Jazz Club", admin)
	system.AddEvent("Art Exhibition", time.Now().Add(48*time.Hour), "Art Gallery", admin)

	fmt.Println("\n--- Guest viewing ---")
	system.ListEvents()

	fmt.Println("\n--- User booking ---")
	system.BookEvent(2, 1, user)

	fmt.Println("\n--- Admin viewing all bookings ---")
	system.ListAllBookings(admin)

	fmt.Println("\n--- User canceling booking ---")
	system.CancelBooking(1, user)

	fmt.Println("\n--- Admin deleting event ---")
	system.DeleteEvent(2, admin)

	fmt.Println("\n--- Final event list ---")
	system.ListEvents()
}
