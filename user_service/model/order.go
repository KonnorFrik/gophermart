package model

import (
	"context"
	"errors"
	"gophermart/model/models"
	"log"
	"sort"

	"github.com/jackc/pgx/v5/pgtype"
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

func NewOrder(order *models.Order, user *models.User) error {
    if order == nil || len(order.Number) == 0 {
        return ErrOrderNothingToCreate
    }

    if dbObj == nil {
        log.Printf("[model.Order/NewOrder]: Lost connection to DB\n")
        connectToPostgres()
        return ErrDataBaseNotConnected
    }

    queries := models.New(dbObj)
    _, err := queries.CreateOrder(context.Background(), models.CreateOrderParams{
        Number: order.Number,
        UserID: pgtype.Int8{Int64: user.ID, Valid: true},
    })

    if err != nil {
        log.Printf("[model.Order/NewOrder]: Error on Create: %q\n", err)
        return err
    }

    // log.Printf("[model.Order/NewOrder]: Create new order#%q, for user(%d)\n", order.Number, user.ID)
    return nil
}

func OrdersRelated(user *models.User) ([]models.Order, error) {
    if user == nil {
        return nil, errors.New("user is nil, can't get orders related to nothing")
    }

    if dbObj == nil {
        log.Printf("[model.Order/OrdersRelated]: Lost connection to DB\n")
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    queries := models.New(dbObj)
    orders, err := queries.UserOrders(context.Background(), pgtype.Int8{Int64: user.ID, Valid: true})

    if err != nil {
        log.Printf("[model.Order/OrdersRelated]: Error on find related: %q\n", err)
        return nil, err
    }

    sort.SliceStable(orders, func(i, j int) bool {
        return orders[j].UploadedAt.Time.Before(orders[i].UploadedAt.Time)
    })

    return orders, nil
}

// TODO: delete order by it id/number 
// TODO: delete all orders belongs to specific user
