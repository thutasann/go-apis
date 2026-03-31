package booking

type MemoryStore struct {
	bookings map[string]Booking // "A2" -> booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (ms *MemoryStore) Book(b Booking) (Booking, error) {
	if _, exists := ms.bookings[b.SeatID]; exists {
		return ms.bookings[b.SeatID], ErrSeatAlreadyBooked
	}

	ms.bookings[b.SeatID] = b
	return ms.bookings[b.SeatID], nil
}

func (ms *MemoryStore) ListBookings(movieID string) []Booking {
	var result []Booking
	for _, b := range ms.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}
