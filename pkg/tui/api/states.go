package api

import "sync"

type Query struct {
	Substring string
	Offset    uint64
	Limit     uint32
}

type Credentials struct {
	Login    string
	Password string
}

type States struct {
	mu            *sync.RWMutex
	query         *Query
	token         string
	credentials   *Credentials
	syncTimestamp string
}

func NewStates() *States {
	return &States{
		mu:          &sync.RWMutex{},
		query:       &Query{},
		credentials: &Credentials{},
	}
}

func (s *States) SetQuery(substring string, offset uint64, limit uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.query.Substring = substring
	s.query.Offset = offset
	s.query.Limit = limit
}

func (s *States) GetQuery() Query {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return *s.query
}

func (s *States) SetTimestamp(timestamp string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.syncTimestamp = timestamp
}

func (s *States) GetTimestamp() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.syncTimestamp
}

func (s *States) SetCredentials(login, password string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.credentials.Login = login
	s.credentials.Password = password
}

func (s *States) GetCredentials() *Credentials {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.credentials
}

func (s *States) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
}

func (s *States) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}
