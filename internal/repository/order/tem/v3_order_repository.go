package tem

import (
	"gocli/internal/model/order/tem"
	"gocli/internal/repository"
	"gocli/util/xlog"
)

// V3OrderRepository 用户分销来源表 模型仓库
type V3OrderRepository struct {
	*repository.DBRepository
}

// NewV3OrderRepository 新用户分销来源表 模型仓库实例
func NewV3OrderRepository(model *tem.V3OrderModel, log *xlog.Log) *V3OrderRepository {
	return &V3OrderRepository{repository.NewDBRepository(model, log)}
}
