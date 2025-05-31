package model

// import (
	// "context"
	// "gophermart/internal/generated/models"
	//    psgs "gophermart/db/postgres"
	// "log"
    // "errors"
// )

// NewUser create a new user in DB.
// Accept raw data from request body.
// func NewUser(user User) (*User, error) {
//     if !user.ValidCredentials() {
//         return nil, ErrUserInvalidData
//     }
//
//     if err := user.HashPassword(); err != nil {
//         return nil, ErrUserInvalidData
//     }
//
//     db := psgs.DB()
//     var err error
//     userDB, err := db.CreateUser(context.TODO(), models.CreateUserParams{
//         Login: user.Login,
//         Password: user.Password,
//     })
//     err = db.WrapError(err)
//
// 	if err != nil {
// 		log.Printf("[model.User/NewUser]: Error on Create: %q\n", err)
// 		return nil, err
// 	}
//
//     toRet := &User{
//         ID: userDB.ID,
//         Login: userDB.Login,
//         // Password: userDB.Password,
//     }
// 	return toRet, nil
// }

// UserByCredentials - Returns a user data stored in DB
// func UserByCredentials(user User) (*User, error) {
//     if !user.ValidCredentials() {
//         return nil, ErrUserInvalidData
//     }
//
//     db := psgs.DB()
//     userDB, err := db.UserByLogin(context.TODO(), user.Login)
//     err = db.WrapError(err)
//
// 	if err != nil {
// 		log.Printf("[model.User/UserByCredentials]: Error: %q\n", err)
// 		return nil, err
// 	}
//
//     if err := user.ComparePassword(userDB.Password); err != nil {
//         log.Printf("[model.User/UserByCredentials]: Error on password compare: %q\n", err)
//         return nil, ErrUserInvalidData
//     }
//
//     toRet := &User{
//         ID: userDB.ID,
//         Login: userDB.Login,
//         // Password: userDB.Password,
//     }
// 	return toRet, nil
// }

// DeleteUserById - Delete user from DB
// func DeleteUserById(id int64) error {
//     db := psgs.DB()
//     err := db.DeleteUser(context.TODO(), id)
//     err = db.WrapError(err)
//
// 	if err != nil {
// 		log.Printf("[model.User/DeleteUser]: Error on delete %q\n", err)
// 		return err
// 	}
//
// 	return nil
// }
