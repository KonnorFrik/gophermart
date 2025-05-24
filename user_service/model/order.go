package model

import (
	"context"
	"errors"
	"gophermart/internal/generated/models"
	"log"
	"sort"
	"time"
)


var (
    ErrOrderAlreadyExist = errors.New("order already exist")
    ErrOrderNothingToCreate = errors.New("nothing to create")
)


type Order struct {
    ID         int64 `json:"-"`
    Accrual    int64 `json:"accrual"`
	Number     string `json:"number"`
	Status     string `json:"status"`
    UploadedAt string `json:"uploaded_at"`
	UserID     int64 `json:"-"`
}

func NewOrder(number string, userID int64) error {
    if dbObj == nil {
        log.Printf("[model.Order/NewOrder]: Lost connection to DB\n")
        connectToPostgres()
        return ErrDataBaseNotConnected
    }

    queries := getQueries()
    defer putQueries(queries)
    _, err := queries.CreateOrder(context.Background(), models.CreateOrderParams{
        Number: number,
        UserID: userID,
    })

    if err != nil {
        log.Printf("[model.Order/NewOrder]: Error on Create: %q\n", err)
        return err
    }

    return nil
}

func OrdersRelated(userID int64) ([]*Order, error) {
    if dbObj == nil {
        log.Printf("[model.Order/OrdersRelated]: Lost connection to DB\n")
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    queries := getQueries()
    defer putQueries(queries)
    orders, err := queries.UserOrders(context.Background(), userID)

    if err != nil {
        log.Printf("[model.Order/OrdersRelated]: Error on find related: %q\n", err)
        return nil, err
    }

    sort.SliceStable(orders, func(i, j int) bool {
        return orders[j].UploadedAt.Time.Before(orders[i].UploadedAt.Time)
    })

    ordersRet := make([]*Order, len(orders))

    for i, v := range orders {
        ordersRet[i] = &Order{
            ID: v.ID,
            Accrual: v.Accrual,
            Number: v.Number,
            Status: string(v.Status),
            UploadedAt: v.UploadedAt.Time.Format(time.RFC3339),
            UserID: v.UserID,
        }
    }

    return ordersRet, nil
}

// TODO: delete order by it id/number 
// TODO: delete all orders belongs to specific user
