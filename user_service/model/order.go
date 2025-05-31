package model

import (
	"context"
	psgs "gophermart/db/postgres"
	"gophermart/internal/generated/models"
	"sort"
	"time"
	"unicode"
)

type Order struct {
    ID         int64 `json:"-"`
    Accrual    int64 `json:"accrual"`
	Number     string `json:"number"`
	Status     string `json:"status"`
    UploadedAt string `json:"uploaded_at"`
	UserID     int64 `json:"-"`
}

// Create - create a new order in db from data in 'o'
// After successfull creatino fill 'o' with created data
func (o *Order) Create(ctx context.Context) error {
    // TODO validate number in another method
    if len(o.Number) == 0 {
        return ErrInvalidData
    }

    for _, r := range o.Number {
        if !unicode.IsDigit(r) {
            return ErrInvalidData
        }
    }
    
    if !validByLUHN(o.Number) {
        return ErrInvalidData
    }


    db := psgs.DB()
    _, err := db.CreateOrder(ctx, models.CreateOrderParams{
        Number: o.Number,
        UserID: o.UserID,
    })

    if err != nil {
        return wrapError(db.WrapError(err))
    }

    return nil
}

// BelongsToUser - returns all orders belongs to one user.
func (o *Order) BelongsToUser(ctx context.Context) ([]Order, error) {
    db := psgs.DB()
    orders, err := db.UserOrders(ctx, o.UserID)

    if err != nil {
        return nil, wrapError(db.WrapError(err))
    }

    // TODO: sort in sql with 'order by'
    sort.SliceStable(orders, func(i, j int) bool {
        return orders[j].UploadedAt.Time.Before(orders[i].UploadedAt.Time)
    })
    ordersRet := make([]Order, len(orders))

    for i, v := range orders {
        ordersRet[i] = Order{
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

func validByLUHN(numbers string) bool {
    return len(numbers) > 0
}

