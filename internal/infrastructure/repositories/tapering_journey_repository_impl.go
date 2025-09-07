// file: internal/infrastructure/repositories/tapering_journey_repository_impl.go
package repositories

import (
	"context"
	"database/sql"

	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"

	"github.com/google/uuid"
)

type TaperingJourneyRepositoryImpl struct {
    db *sql.DB
}

func NewTaperingJourneyRepositoryImpl(db *sql.DB) repositories.TaperingJourneyRepository {
    return &TaperingJourneyRepositoryImpl{db: db}
}

// Create inserts a new tapering journey record into the database.
func (r *TaperingJourneyRepositoryImpl) Create(ctx context.Context, journey *entities.TaperingJourney) error {
    query := `
        INSERT INTO tapering_journeys 
        (id, user_id, goal_frequency, start_limit, target_limit, reduction_step, step_period, current_limit, status_id, started_at, completed_at, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW())`
    _, err := r.db.ExecContext(ctx, query, journey.ID, journey.UserID, journey.GoalFrequency, journey.StartLimit, journey.TargetLimit, journey.ReductionStep, journey.StepPeriod, journey.CurrentLimit, journey.StatusID, journey.StartedAt, journey.CompletedAt)
    if err != nil {
        return err
    }
    return nil
}

// Update modifies an existing tapering journey record.
func (r *TaperingJourneyRepositoryImpl) Update(ctx context.Context, journey *entities.TaperingJourney) error {
    query := `
        UPDATE tapering_journeys 
        SET current_limit=$1, status_id=$2, completed_at=$3, updated_at=NOW()
        WHERE id=$4 AND user_id=$5`
    _, err := r.db.ExecContext(ctx, query, journey.CurrentLimit, journey.StatusID, journey.CompletedAt, journey.ID, journey.UserID)
    if err != nil {
        return err
    }
    return nil
}

// GetByUserID returns all tapering journeys for a given user.
func (r *TaperingJourneyRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.TaperingJourney, error) {
    query := `SELECT id, user_id, goal_frequency, start_limit, target_limit, reduction_step, step_period, current_limit, status_id, started_at, completed_at, created_at, updated_at 
              FROM tapering_journeys WHERE user_id=$1 ORDER BY created_at DESC`
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var journeys []*entities.TaperingJourney
    for rows.Next() {
        j := &entities.TaperingJourney{}
        err = rows.Scan(&j.ID, &j.UserID, &j.GoalFrequency, &j.StartLimit, &j.TargetLimit, &j.ReductionStep, &j.StepPeriod, &j.CurrentLimit, &j.StatusID, &j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
        if err != nil {
            return nil, err
        }
        journeys = append(journeys, j)
    }
    return journeys, nil
}

// GetByID returns a tapering journey by its ID.
func (r *TaperingJourneyRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.TaperingJourney, error) {
    j := &entities.TaperingJourney{}
    query := `SELECT id, user_id, goal_frequency, start_limit, target_limit, reduction_step, step_period, current_limit, status_id, started_at, completed_at, created_at, updated_at 
              FROM tapering_journeys WHERE id=$1`
    err := r.db.QueryRowContext(ctx, query, id).Scan(&j.ID, &j.UserID, &j.GoalFrequency, &j.StartLimit, &j.TargetLimit, &j.ReductionStep, &j.StepPeriod, &j.CurrentLimit, &j.StatusID, &j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return j, nil
}

// Delete removes a tapering journey by ID and userID.
func (r *TaperingJourneyRepositoryImpl) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
    query := `DELETE FROM tapering_journeys WHERE id=$1 AND user_id=$2`
    _, err := r.db.ExecContext(ctx, query, id, userID)
    return err
}
