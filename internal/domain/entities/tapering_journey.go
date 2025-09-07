// file: internal/domain/entities/tapering_journey.go
package entities

import (
	"time"

	"github.com/google/uuid"
)

// TaperingJourney represents a user's coffee consumption reduction plan over time.
// It tracks the initial goal, target goal, reduction steps, timing, and current status.
// The goal can be defined with different frequencies (daily, weekly, monthly).
// The struct supports flexible tapering by adjusting limits in defined intervals.
type TaperingJourney struct {
    // ID is the unique identifier for this tapering journey.
    ID uuid.UUID `json:"id" db:"id"`

    // UserID references the user this tapering journey belongs to.
    UserID uuid.UUID `json:"user_id" db:"user_id"`

    // GoalFrequency indicates the frequency unit for the consumption goal.
    // E.g., 1 = daily, 2 = weekly, 3 = monthly.
    GoalFrequency GoalFrequency `json:"goal_frequency" db:"goal_frequency"`

    // StartLimit is the initial consumption limit (e.g., cups) set at the journey start,
    // expressed in units per the GoalFrequency (per day, per week, etc.).
    StartLimit int `json:"start_limit" db:"start_limit"`

    // TargetLimit is the desired final consumption limit to reach by the end of the journey.
    TargetLimit int `json:"target_limit" db:"target_limit"`

    // ReductionStep is the amount by which the consumption limit decreases at each step.
    ReductionStep int `json:"reduction_step" db:"reduction_step"`

    // StepPeriod specifies the number of days between each reduction step.
    // It works in conjunction with GoalFrequency (e.g., StepPeriod=7 means reduce every 7 days).
    StepPeriod int `json:"step_period" db:"step_period"`

    // CurrentLimit shows the current active limit at this point in the tapering journey.
    CurrentLimit int `json:"current_limit" db:"current_limit"`

    // StatusID represents the current status of the tapering journey.
    // Possible values include Active, Paused, and Completed.
    StatusID GoalStatus `json:"status_id" db:"status_id"`

    // StartedAt is the timestamp when the tapering journey was initiated.
    StartedAt time.Time `json:"started_at" db:"started_at"`

    // CompletedAt is an optional timestamp indicating when the journey was completed.
    CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`

    // CreatedAt is the timestamp when the journey record was created.
    CreatedAt time.Time `json:"created_at" db:"created_at"`

    // UpdatedAt is the timestamp of the last update to the journey record.
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GoalFrequency defines the frequency unit for tapering goals.
type GoalFrequency int

const (
    // Daily frequency for goals counted per day.
    Daily GoalFrequency = iota + 1

    // Weekly frequency for goals counted per week.
    Weekly

    // Monthly frequency for goals counted per month.
    Monthly
)

// GoalStatus represents the status of the tapering journey.
type GoalStatus int

const (
    // Active indicates the journey is currently in progress.
    Active GoalStatus = iota + 1

    // Paused indicates the journey is temporarily halted.
    Paused

    // Completed indicates the journey has been finished successfully.
    Completed

	// User cancelled the journey.
	Cancelled
)