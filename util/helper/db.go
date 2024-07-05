package helper

import (
	"fmt"
	"strings"
)

func GetDatabaseName(dsn string) (string, error) {
	// 从DSN中获取数据库信息部分
	parts := strings.Split(dsn, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid DSN: missing database name")
	}

	// 数据库信息部分
	dbInfo := parts[1]

	// 去掉参数部分
	dbName := strings.Split(dbInfo, "?")[0]
	if dbName == "" {
		return "", fmt.Errorf("invalid DSN: missing database name")
	}

	return dbName, nil
}
