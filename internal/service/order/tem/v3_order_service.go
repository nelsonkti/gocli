package tem

import (
	"gocli/internal/server"
)

// V3OrderService 用户分销来源表 模型服务
type V3OrderService struct {
	svc *server.SvcContext
}

// NewV3OrderService 新用户分销来源表 模型服务实例
func NewV3OrderService(svc *server.SvcContext) *V3OrderService {
	return &V3OrderService{svc: svc}
}
