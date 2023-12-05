package service

func (s *Service) Ping() error {
	return s.store.Ping()
}
