package model

import (
	"gophermart/internal/generated/models"
	"sync"
)

var _queryPool = sync.Pool{
    New: func() any { return models.New(dbObj) },
}

func getQueries() *models.Queries {
    return _queryPool.Get().(*models.Queries)
}

func putQueries(obj *models.Queries) {
    _queryPool.Put(obj)
}
