package test

import (
	"gocli/internal/model/test"
	"gocli/internal/repository"
	"gocli/util/xlog"
)

// V3OrderRepository 测试仓库
type V3OrderRepository struct {
	*repository.DBRepository
}

// NewV3OrderRepository 新测试仓库实例
func NewV3OrderRepository(model *test.V3OrderModel, log *xlog.Log) *V3OrderRepository {
	return &V3OrderRepository{repository.NewDBRepository(model, log)}
}
