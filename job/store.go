package job

type Store struct {
	storage Storage
}

type Storage interface {
	Get(key []byte) (*Job, error)
	Put(key []byte, j *Job) error
}

func NewStore(s Storage) *Store {
	return &Store{storage: s}
}

func (s *Store) CreateJob(url string) (*Job, error) {
	j := &Job{URL: url}
	n, err := newNamer(j)
	if err != nil {
		return nil, err
	}
	for {
		nm := n.name()

		v, err := s.storage.Get([]byte(nm))
		if err != nil {
			return nil, err
		}

		// No existing entry
		if v == nil {
			j.Name = nm
			if err := s.storage.Put([]byte(j.Name), j); err != nil {
				return nil, err
			}
			return j, err
		}

		// Already exist
		if v.URL == j.URL {
			return nil, nil // No need to create a job
		}

		// Name conflict with different URL
		n.next() // Generate next name and retry
	}
}

func (s *Store) Put(j *Job) error {
	return s.storage.Put([]byte(j.Name), j)
}
