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
    ErrOrderInvalidNumber = errors.New("invalid order number")
)

func NewOrder(number string, userID int64) error {
    if dbObj == nil {
        connectToPostgres()
        return ErrDataBaseNotConnected
    }

    if len(number) == 0 {
        return ErrOrderInvalidNumber
    }

    for _, r := range number {
        if !unicode.IsDigit(r) {
            return ErrOrderInvalidNumber
        }
    }
    
    if !validByLUHN(number) {
        return ErrOrderInvalidNumber
    }

    queries := getQueries()
    defer putQueries(queries)
    _, err := queries.CreateOrder(context.TODO(), models.CreateOrderParams{
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
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    queries := getQueries()
    defer putQueries(queries)
    orders, err := queries.UserOrders(context.TODO(), userID)
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return nil, err
    }

    if err != nil {
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
