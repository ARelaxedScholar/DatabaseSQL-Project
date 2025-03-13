package models

import (
	"errors"
	"time"
)

type Problem struct {
	id             int
	severity       ProblemSeverity
	description    string
	signaledWhen   time.Time
	isResolved     bool
	resolutionDate time.Time // null only if unresolved
}

func (self Problem) validate() error {
	var err error
	switch {
	case self.id < 0:
		err = errors.New("Problem's id cannot be negative.")
	case !self.severity.isValid():
		err = errors.New("Problem contains invalid severity variant.")
	case self.description == "":
		err = errors.New("Problem's description cannot be empty.")
	case self.resolutionDate.IsZero() && self.isResolved:
		err = errors.New("Problem is resolved, but no resolution time was passed.")
	}
	return err // note that this will be nil if none of the cases were hit

}
