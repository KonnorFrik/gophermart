package model

// import (
// 	"context"
// 	"errors"
// 	psgs "gophermart/db/postgres"
//     "gophermart/internal/generated/models"
// 	"sort"
// 	"time"
// 	"unicode"
// )

// func NewOrder(number string, userID int64) error {
//     if len(number) == 0 {
//         return ErrOrderInvalidNumber
//     }
//
//     for _, r := range number {
//         if !unicode.IsDigit(r) {
//             return ErrOrderInvalidNumber
//         }
//     }
//
//     if !validByLUHN(number) {
//         return ErrOrderInvalidNumber
//     }
//
//
//     db := psgs.DB()
//     _, err := db.CreateOrder(context.TODO(), models.CreateOrderParams{
//         Number: number,
//         UserID: userID,
//     })
//     err = db.WrapError(err)
//
//     return err
// }

// func OrdersRelated(userID int64) ([]*Order, error) {
//     db := psgs.DB()
//     orders, err := db.UserOrders(context.TODO(), userID)
//     err = db.WrapError(err)
//
//     if err != nil {
//         return nil, err
//     }
//
//     // TODO: sort in sql with 'order by'
//     sort.SliceStable(orders, func(i, j int) bool {
//         return orders[j].UploadedAt.Time.Before(orders[i].UploadedAt.Time)
//     })
//     ordersRet := make([]*Order, len(orders))
//
//     for i, v := range orders {
//         ordersRet[i] = &Order{
//             ID: v.ID,
//             Accrual: v.Accrual,
//             Number: v.Number,
//             Status: string(v.Status),
//             UploadedAt: v.UploadedAt.Time.Format(time.RFC3339),
//             UserID: v.UserID,
//         }
//     }
//
//     return ordersRet, nil
// }
