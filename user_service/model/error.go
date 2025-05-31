/*
File for predefined error
Maps errors from postgres/db pkg into model errors
*/
package model

import (
	"errors"
	"fmt"
	psgs "gophermart/db/postgres"
	"log"
)

var (
    // ErrDoesNotExist - data not exist in db
	ErrDoesNotExist = errors.New("does not exist")
    // ErrAlreadyExist - data already exist in db
    ErrAlreadyExist = errors.New("already exist")
    // ErrInvalidData - data is corrupted, missed, wrong format, etc
    ErrInvalidData = errors.New("invalid data")
    // ErrDbNoAccess - no access to database
    ErrDbNoAccess = psgs.ErrDataBaseNotConnected
    // ErrUnknown - any other not documented error
    ErrUnknown = psgs.ErrUnknown
)

func wrapError(err error) error {
    if err == nil {
        return nil
    }

    switch {
    case errors.Is(err, psgs.ErrInvalidConfig):
        return fmt.Errorf("%w: %w", ErrUnknown, err)
    case errors.Is(err, psgs.ErrDataBaseNotConnected):
        return fmt.Errorf("%w: %w", ErrDbNoAccess, err)
    case errors.Is(err, psgs.ErrConstraintUniqueViolation):
        return fmt.Errorf("%w: %w", ErrAlreadyExist, err)
    case errors.Is(err, psgs.ErrConstraintForeignKeyViolation):
        return fmt.Errorf("%w: %w", ErrInvalidData, err)
    case errors.Is(err, psgs.ErrConstraintCheckViolation):
        return fmt.Errorf("%w: %w", ErrInvalidData, err)
    case errors.Is(err, psgs.ErrUnknown):
        return fmt.Errorf("%w: %w", ErrUnknown, err)
    }

    log.Printf("[model.wrapError]: CAN'T MAP ERR: %q\n", err)
    return fmt.Errorf("%w: %w", ErrUnknown, err)
}
