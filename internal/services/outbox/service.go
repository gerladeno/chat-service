package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	jobsrepo "github.com/gerladeno/chat-service/internal/repositories/jobs"
	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/types"
)

const (
	serviceName                       = "outbox"
	reasonJobNotFound                 = "not_found"
	reasonFailedAttemptsLimitExceeded = "too_many_errors"
)

type jobsRepository interface {
	CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
	FindAndReserveJob(ctx context.Context, until time.Time) (jobsrepo.Job, error)
	CreateFailedJob(ctx context.Context, name, payload, reason string) error
	DeleteJob(ctx context.Context, jobID types.JobID) error
}

type transactor interface {
	RunInTx(context.Context, func(ctx context.Context) error) error
}

var ErrJobAlreadyExists = errors.New("job already exists")

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	workers    int             `option:"mandatory" validate:"min=1,max=32"`
	idleTime   time.Duration   `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration   `option:"mandatory" validate:"min=1s,max=10m"`
	jobsRepo   *jobsrepo.Repo  `option:"mandatory"`
	db         *store.Database `option:"mandatory"`
}

type Service struct {
	workers    int
	idleTime   time.Duration
	reserveFor time.Duration
	registry   map[string]Job
	jobsRepo   jobsRepository
	db         transactor
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating outbox service options: %v", err)
	}
	return &Service{
		registry:   make(map[string]Job),
		workers:    opts.workers,
		idleTime:   opts.idleTime,
		reserveFor: opts.reserveFor,
		jobsRepo:   opts.jobsRepo,
		db:         opts.db,
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	if _, ok := s.registry[job.Name()]; ok {
		return ErrJobAlreadyExists
	}
	s.registry[job.Name()] = job
	return nil
}

func (s *Service) MustRegisterJob(job Job) {
	if err := s.RegisterJob(job); err != nil {
		panic(err)
	}
}

func (s *Service) Run(ctx context.Context) error {
	for i := 0; i < s.workers; i++ {
		go s.startWorker(ctx, i)
	}
	return nil
}

func (s *Service) startWorker(ctx context.Context, workerID int) {
	log := zap.L().With(zap.String("service", serviceName), zap.Int("worker_id", workerID))
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		err = s.execute(ctx, log)
		switch {
		case errors.Is(err, jobsrepo.ErrNoJobs):
			log.Info(fmt.Sprintf("out of jobs, idling for %d milliseconds", s.idleTime.Milliseconds()))
			time.Sleep(s.idleTime)
		case err != nil:
			log.With(zap.Error(err)).Warn("execution failed, proceeding")
		}
	}
}

func (s *Service) execute(ctx context.Context, log *zap.Logger) (err error) {
	task, err := s.jobsRepo.FindAndReserveJob(ctx, time.Now().Add(s.reserveFor))
	if err != nil {
		return fmt.Errorf("get a new task: %w", err)
	}
	l := log.With(
		zap.String("job", task.Name),
		zap.String("payload", task.Payload),
		zap.Int("attempts", task.Attempts))
	l.Info("executing task")
	defer func() {
		if err == nil {
			l.Info("success")
		} else {
			l.Warn("failed")
		}
	}()

	job, ok := s.registry[task.Name]
	if !ok {
		return s.moveToDLQ(ctx, task, reasonJobNotFound)
	}

	ctx, cancel := context.WithTimeout(ctx, job.ExecutionTimeout())
	defer cancel()
	if err = job.Handle(ctx, task.Payload); err != nil {
		if task.Attempts >= job.MaxAttempts() {
			go func() {
				if err := s.moveToDLQ(context.Background(), task, reasonFailedAttemptsLimitExceeded); err != nil {
					l.Warn("err during handling an error", zap.Error(err))
				}
			}()
		}
		return fmt.Errorf("handling a job %v: %v", task, err)
	}
	if err = s.jobsRepo.DeleteJob(ctx, task.ID); err != nil {
		return fmt.Errorf("delete not found job: %v", err)
	}
	return nil
}

func (s *Service) moveToDLQ(ctx context.Context, task jobsrepo.Job, reason string) error {
	return s.db.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.jobsRepo.DeleteJob(ctx, task.ID); err != nil {
			return fmt.Errorf("delete not found job: %v", err)
		}
		if err := s.jobsRepo.CreateFailedJob(ctx, task.Name, task.Payload, reason); err != nil {
			return fmt.Errorf("fail not found job: %v", err)
		}
		return nil
	})
}
