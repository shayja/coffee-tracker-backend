// file: internal/infrastructure/http/dto/tapering_journey_dto.go
package dto

import "time"

type CreateTaperingJourneyRequest struct {
	  // GoalFrequency defines the goal interval, e.g., daily=1, weekly=2, monthly=3.
    GoalFrequency int `json:"goal_frequency" binding:"required"`

    // StartLimit is the initial consumption limit (e.g., cups) per frequency unit.
    StartLimit int `json:"start_limit" binding:"required,min=0"`

    // TargetLimit is the desired consumption limit to reach.
    TargetLimit int `json:"target_limit" binding:"required,min=0"`

    // ReductionStep specifies how much to reduce consumption at each step.
    ReductionStep int `json:"reduction_step" binding:"required,min=1"`

    // StepPeriod is the number of days between each reduction step.
    StepPeriod int `json:"step_period" binding:"required,min=1"`

    // StartedAt is optional. If not provided, the backend sets it to current time.
    StartedAt *time.Time `json:"started_at,omitempty"`
}
