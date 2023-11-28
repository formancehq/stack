package service

func (s *Service) Ping() error {
	return newStorageError(s.store.Ping(), "ping")
}
