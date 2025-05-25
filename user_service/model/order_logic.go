package model

import (
	"context"
	"errors"
	"gophermart/internal/generated/models"
	"sort"
	"time"
	"unicode"
)

var (
    // ErrOrderAlreadyExist = errors.New("order already exist")
    ErrOrderInvalidInput = errors.New("invalid order number")
)

func NewOrder(number string, userID int64) error {
    if dbObj == nil {
        // log.Printf("[model.Order/NewOrder]: Lost connection to DB\n")
        connectToPostgres()
        return ErrDataBaseNotConnected
    }

    if len(number) == 0 {
        // log.Printf("[modes.Order/NewOrder]: Zero data length\n")
        return ErrOrderInvalidInput
    }

    for _, r := range number {
        if !unicode.IsDigit(r) {
            // log.Printf("[modes.Order/NewOrder]: Invalid input data: %q\n", number)
            return ErrOrderInvalidInput
        }
    }
    
    if !validByLUHN(number) {
        // log.Printf("[modes.Order/NewOrder]: Invalid by LUHN input: %q\n", number)
        return ErrOrderInvalidInput
    }

    queries := getQueries()
    defer putQueries(queries)
    _, err := queries.CreateOrder(context.Background(), models.CreateOrderParams{
        Number: number,
        UserID: userID,
    })
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return err
    }

    return err
}

func OrdersRelated(userID int64) ([]*Order, error) {
    if dbObj == nil {
        // log.Printf("[model.Order/OrdersRelated]: Lost connection to DB\n")
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    queries := getQueries()
    defer putQueries(queries)
    orders, err := queries.UserOrders(context.Background(), userID)
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return nil, err
    }

    if err != nil {
        // log.Printf("[model.Order/OrdersRelated]: Error on find: %q\n", err)
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
