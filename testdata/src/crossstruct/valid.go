package crossstruct

type Service struct {
	repo *Repo
}

func (s *Service) serve() {
	s.repo.find()
}

type Repo struct{}

func (r *Repo) find() {}
