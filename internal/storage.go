package internal

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(event CreateEventRequest) {

}
