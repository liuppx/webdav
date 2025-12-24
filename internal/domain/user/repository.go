package user

import "context"

// Repository 用户仓储接口
type Repository interface {
	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*User, error)
	
	// FindByWalletAddress 根据钱包地址查找用户
	FindByWalletAddress(ctx context.Context, address string) (*User, error)
	
	// Save 保存用户
	Save(ctx context.Context, user *User) error
	
	// Delete 删除用户
	Delete(ctx context.Context, username string) error
	
	// List 列出所有用户
	List(ctx context.Context) ([]*User, error)
}

