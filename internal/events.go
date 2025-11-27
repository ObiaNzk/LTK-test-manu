package internal

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HelloWorld() string {
	return "Hello, World!"
}
