package models

import (
	"errors"
	"time"
)

type Problem struct {
	ID             int
	Severity       ProblemSeverity
	Description    string
	SignaledWhen   time.Time
	IsResolved     bool
	ResolutionDate time.Time // null only if unresolved
}

func (self Problem) validate() error {
	var err error
	switch {
	case self.ID < 0:
		err = errors.New("Problem's id cannot be negative.")
	case !self.Severity.isValid():
		err = errors.New("Problem contains invalid severity variant.")
	case self.Description == "":
		err = errors.New("Problem's description cannot be empty.")
	case self.ResolutionDate.IsZero() && self.IsResolved:
		err = errors.New("Problem is resolved, but no resolution time was passed.")
	}
	return err // note that this will be nil if none of the cases were hit

}
