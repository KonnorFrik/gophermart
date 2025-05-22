package model

import (
	"errors"
	"log"
	"sort"
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
    ID         uint      `gorm:"primarykey" json:"-"`
    Accrual    uint      `json:"accrual"`
    Number     string    `gorm:"unique" json:"number"`
    Status     Status    `gorm:"type:bigint" json:"status"`
    UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
    UserID     uint      `json:"-"`
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

func OrdersRelated(user *User) ([]Order, error) {
    if user == nil {
        return nil, errors.New("user is nil, can't get orders related to nothing")
    }

    if dbObj == nil {
        log.Printf("[model.Order/OrdersRelated]: Lost connection to DB\n")
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    var orders []Order
    err := dbObj.Model(user).Association("Orders").Find(&orders)

    if err != nil {
        log.Printf("[model.Order/OrdersRelated]: Error on find related: %q\n", err)
        return nil, err
    }

    sort.SliceStable(orders, func(i, j int) bool {
        return orders[j].UploadedAt.Before(orders[i].UploadedAt)
    })

    return orders, nil
}
