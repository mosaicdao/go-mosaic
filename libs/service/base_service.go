package service

import (
	"errors"
	"sync"

	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("service")

	// ErrAlreadyRunning is returned on a call to run an already running service.
	ErrAlreadyRunning = errors.New("already running")
)

type BaseService struct {
	name      string
	isRunning bool
	quit      chan struct{}
	mux       sync.Mutex
	impl      Servicable
}

// NewBaseService creates a new service from a servicable object.
func NewBaseService(name string, impl Servicable) *BaseService {
	return &BaseService{
		name:      name,
		isRunning: false,
		quit:      make(chan struct{}),
		impl:      impl,
	}
}

// Start the service.
func (s *BaseService) Start() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.isRunning {
		log.Errorf("Not starting %v service because already running (impl: %v)", s.name, s.impl)
		return ErrAlreadyRunning
	}

	if err := s.impl.OnStart(); err != nil {
		log.Errorf("Not starting %v service: %w (impl: %v)", s.name, err, s.impl)
		return err
	}
	s.quit = make(chan struct{})
	s.isRunning = true
	log.Infof("Starting service %v (impl: %v).", s.name, s.impl)
	return nil
}

// Stop the service.
func (s *BaseService) Stop() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if !s.isRunning {
		return
	}

	s.impl.OnStop()
	close(s.quit)
	s.isRunning = false
}

// Wait blocks until the service is stopped.
func (s *BaseService) Wait() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if !s.isRunning {
		return
	}
	s.mux.Unlock()

	<-s.quit
}

// IsRunning returns true when the service is running.
func (s *BaseService) IsRunning() bool {
	return s.isRunning
}

// String returns a string representation of the service.
func (s *BaseService) String() string {
	return s.name
}
