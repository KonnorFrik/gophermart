/*
Connection manager
Package for store database object as singletone
Has auto reconnection, connection retries will be until connect is success
*/
package postgres

import (
	// "log"
	// "sync"
	// "sync/atomic"
	// "time"
    "errors"
    "fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
    // DbDefaultSslmode = "disable"
    // DbDefaultTimezone = "Asia/Yekaterinburg"

    // reconnectDelay = time.Minute
)

// const (
//     stateNotConnected int32 = iota
//     stateConnected
// )

var (
    // dbSingletone DB
    // dbConfig DbConfig

    // reconnectRequest chan struct{}
    // stopReconnect chan struct{}

    // state int32 

    ErrInvalidConfig = errors.New("invalid database configuration")
)

// func init() {
//     reconnectRequest = make(chan struct{}, 1)
//     stopReconnect = make(chan struct{})
// }

// type DB struct {
//     sync.Mutex
//     DB *gorm.DB
// }

// SaveConfig - save configuration for connect to database.
// It will be used by connection retrier
// func SaveConfig(config DbConfig) error {
//     if !config.isValid() {
//         return ErrInvalidConfig
//     }
//
//     dbConfig = config
//     log.Printf("[PostgreSQL/SaveConfig]: config saved: %q\n", config.String())
//     return nil
// }

// Connect - connect to db with saved config
func Connect(config DbConfig) (*gorm.DB, error) {
    // select {
    // case reconnectRequest <- struct{}{}:
    // default:
    // }

    if !config.isValid() {
        return nil, ErrInvalidConfig
    }

    var db *gorm.DB
    var err error
    connectionInfo := config.String()
    db, err = gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})

    if err != nil {
        return nil, err
    }

    return db, nil
}

// func Get() *DB {
//     dbSingletone.Lock()
//     defer dbSingletone.Unlock()
//
//     if dbSingletone.DB == nil {
//         if atomic.CompareAndSwapInt32(&state, stateNotConnected, stateConnected) {
//             log.Printf("[PostgreSQL/Get]: DB is nil, switch state and run retrier\n")
//             go connectionRetrier(reconnectRequest, stopReconnect)
//
//         } else {
//             log.Printf("[PostgreSQL/Get]: DB is nil, send reconnect request\n")
//             select {
//             case reconnectRequest <- struct{}{}:
//             default:
//             }
//         }
//     }
//
//     return &dbSingletone
// }

// func CloseConnection() {
//     select {
//     case stopReconnect <- struct{}{}:
//     default:
//     }
//
//     dbSingletone.Lock()
//     defer dbSingletone.Unlock()
//
//     if dbSingletone.DB != nil {
//         if atomic.CompareAndSwapInt32(&state, stateConnected, stateNotConnected) {
//             sql, err := dbSingletone.DB.DB()
//
//             if err == nil {
//                 sql.Close()
//                 dbSingletone.DB = nil
//                 log.Printf("[PostgreSQL/CloseConnection]: Close DB connection\n")
//
//             } else {
//                 log.Printf("[PostgreSQL/CloseConnection]: Error on get DB: %q\n", err)
//             }
//         }
//     }
// }

// func connectionRetrier(reconnect, shutdown <-chan struct{}) {
//     ticker := time.NewTicker(reconnectDelay)
//     var err error
//
//     for {
//         select {
//         case <-ticker.C:
//             dbSingletone.Lock()
//
//             if atomic.LoadInt32(&state) == stateConnected && dbSingletone.DB == nil {
//                 log.Printf("[PostgreSQL/retrier]: DB is not connected. Try to connect\n")
//                 dbSingletone.DB, err = gorm.Open(postgres.Open(dbConfig.String()), &gorm.Config{})
//
//                 if err != nil {
//                     log.Printf("[PostgreSQL/retrier]: Erron on reconnect: %q\n", err)
//                 }
//             }
//
//             dbSingletone.Unlock()
//
//         case <-reconnect:
//             dbSingletone.Lock()
//             log.Printf("[PostgreSQL/retrier]: Got reconnect request\n")
//             dbSingletone.DB, err = gorm.Open(postgres.Open(dbConfig.String()), &gorm.Config{})
//
//             if err != nil {
//                 log.Printf("[PostgreSQL/retrier]: Erron on reconnect: %q\n", err)
//             }
//
//             dbSingletone.Unlock()
//
//         // kill goroutine
//         case <-shutdown:
//             log.Printf("[PostgreSQL/retrier]: shutdown\n")
//             return
//         }
//     }
// }

type DbConfig struct {
    Host string
    User string 
    Password string 
    DbName string 
    Port string
    SSLMode bool
    TimeZone string
}

func (dc DbConfig) String() string {
    sslMode := "disable"

    if dc.SSLMode {
        sslMode = "enable"
    }

    return fmt.Sprintf(
        "host=%s database=%s user=%s password=%s port=%s sslmode=%s TimeZone=%s",
        dc.Host,
        dc.DbName,
        dc.User,
        dc.Password,
        dc.Port,
        sslMode,
        dc.TimeZone,
    )
}

func (dc DbConfig) isValid() bool {
    return dc.Host != "" && dc.User != "" && dc.Password != "" && dc.Port != "" && dc.TimeZone != "" 
}
