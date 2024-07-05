package tem

import (
	"gocli/internal/model/order/tem"
	"gocli/internal/repository"
	"gocli/util/xlog"
)

// CodeUserFromMaRepository 用户分销来源表 模型仓库
type CodeUserFromMaRepository struct {
	*repository.DBRepository
}

// NewCodeUserFromMaRepository 新用户分销来源表 模型仓库实例
func NewCodeUserFromMaRepository(model *tem.CodeUserFromMaModel, log *xlog.Log) *CodeUserFromMaRepository {
	return &CodeUserFromMaRepository{repository.NewDBRepository(model, log)}
}
