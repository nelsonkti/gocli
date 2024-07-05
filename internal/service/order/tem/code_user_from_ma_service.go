package tem

import (
	"gocli/internal/server"
)

// CodeUserFromMaService 用户分销来源表 模型服务
type CodeUserFromMaService struct {
	svc *server.SvcContext
}

// NewCodeUserFromMaService 新用户分销来源表 模型服务实例
func NewCodeUserFromMaService(svc *server.SvcContext) *CodeUserFromMaService {
	return &CodeUserFromMaService{svc: svc}
}
