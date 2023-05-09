package jobsrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/store/job"
	"github.com/gerladeno/chat-service/internal/types"
)

var ErrNoJobs = errors.New("no jobs found")

type Job struct {
	ID       types.JobID
	Name     string
	Payload  string
	Attempts int
}

func (r *Repo) FindAndReserveJob(ctx context.Context, until time.Time) (Job, error) {
	var result Job
	err := r.db.RunInTx(ctx, func(ctx context.Context) error {
		foundJob, err := r.db.Job(ctx).Query().Where(job.And(
			job.AvailableAtLTE(time.Now()),
			job.ReservedUntilLTE(time.Now()),
		)).Order(job.ByCreatedAt()).ForUpdate().First(ctx)
		switch {
		case store.IsNotFound(err):
			return ErrNoJobs
		case err != nil:
			return fmt.Errorf("finding a job: %v", err)
		}
		foundJob, err = r.db.Job(ctx).UpdateOne(foundJob).SetReservedUntil(until).AddAttempts(1).Save(ctx)
		if err != nil {
			return fmt.Errorf("updating found job: %v", err)
		}
		result = Job{
			ID:       foundJob.ID,
			Name:     foundJob.Name,
			Payload:  foundJob.Payload,
			Attempts: foundJob.Attempts,
		}
		return nil
	})
	return result, err
}

func (r *Repo) CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	newJob, err := r.db.Job(ctx).Create().
		SetName(name).
		SetPayload(payload).
		SetAvailableAt(availableAt).
		Save(ctx)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("creating a job: %v", err)
	}
	return newJob.ID, nil
}

func (r *Repo) CreateFailedJob(ctx context.Context, name, payload, reason string) error {
	_, err := r.db.FailedJob(ctx).Create().
		SetName(name).
		SetPayload(payload).
		SetReason(reason).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating a job: %v", err)
	}
	return nil
}

func (r *Repo) DeleteJob(ctx context.Context, jobID types.JobID) error {
	err := r.db.Job(ctx).DeleteOneID(jobID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting a job: %v", err)
	}
	return nil
}
