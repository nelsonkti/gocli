package test

import (
	"gocli/internal/model/internal/model/test"
	"gocli/internal/repository"
	"gocli/util/xlog"
)

// CodeUserFromMaRepository 测试仓库
type CodeUserFromMaRepository struct {
	*repository.DBRepository
}

// NewCodeUserFromMaRepository 新测试仓库实例
func NewCodeUserFromMaRepository(model *test.CodeUserFromMaModel, log *xlog.Log) *CodeUserFromMaRepository {
	return &CodeUserFromMaRepository{repository.NewDBRepository(model, log)}
}
