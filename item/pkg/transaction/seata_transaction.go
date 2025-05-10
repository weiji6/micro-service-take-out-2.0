package transaction

import (
	"context"

	"github.com/seata/seata-go/pkg/client"
	"github.com/seata/seata-go/pkg/tm"
)

// SeataTransaction Seata分布式事务管理器
type SeataTransaction struct {
	client *client.Client
}

// NewSeataTransaction 创建新的Seata事务管理器
func NewSeataTransaction(client *client.Client) *SeataTransaction {
	return &SeataTransaction{
		client: client,
	}
}

// WithTransaction 在事务中执行操作
func (t *SeataTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name: "business-transaction",
	}, fn)
}

// Begin 开始事务
func (t *SeataTransaction) Begin(ctx context.Context) (context.Context, error) {
	return tm.Begin(ctx, "business-transaction")
}

// Commit 提交事务
func (t *SeataTransaction) Commit(ctx context.Context) error {
	return tm.Commit(ctx)
}

// Rollback 回滚事务
func (t *SeataTransaction) Rollback(ctx context.Context) error {
	return tm.Rollback(ctx)
} 