package model

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

const (
	statusNew Status = iota
	statusProcessing
	statusInvalid
	statusProcessed
)

var (
    ErrOrderAlreadyExist = errors.New("order already exist")
    ErrOrderNothingToCreate = errors.New("nothing to create")
)

type Status int

type Order struct {
    ID         uint      `gorm:"primarykey"`
	Accrual    uint
    Number     string    `gorm:"unique"`
	Status     Status    `gorm:"type:bigint"`
	UploadedAt time.Time `gorm:"autoCreateTime"`
    UserID     uint
}

func (s Status) String() string {
    switch s {
    case statusNew:
        return "NEW"
    case statusProcessing:
        return "PROCESSING"
    case statusInvalid:
        return "INVALID"
    case statusProcessed:
        return "PROCESSED"

    default:
        return ""
    }
}

func NewOrder(order *Order, user *User) error {
    if order == nil || len(order.Number) == 0 {
        return ErrOrderNothingToCreate
    }

    if dbObj == nil {
        log.Printf("[model.Order/NewOrder]: Lost connection to DB\n")
        connectToPostgres()
        return ErrDataBaseNotConnected
    }

    // log.Printf("[model.Order/NewOrder]: Create a new order(%+v) for user(%+v)\n", order, user)
    var err error
    err = dbObj.Model(user).Association("Orders").Append(order)

    if errors.Is(err, gorm.ErrDuplicatedKey) {
        log.Printf("[model.Order/NewOrder]: Order already exist: %q\n", order.Number)
        return ErrOrderAlreadyExist
    }

    if err != nil {
        log.Printf("[model.Order/NewOrder]: Error on Create: %q\n", err)
        return err
    }

    return nil
}
