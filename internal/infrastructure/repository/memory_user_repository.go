package repository

import (
	"context"
	"strings"
	"sync"
	
	"github.com/yeying-community/webdav/internal/domain/user"
	"github.com/yeying-community/webdav/internal/infrastructure/config"
	"github.com/yeying-community/webdav/internal/infrastructure/crypto"
)

// MemoryUserRepository 内存用户仓储
type MemoryUserRepository struct {
	users           map[string]*user.User // username -> user
	walletAddresses map[string]*user.User // wallet_address -> user
	mu              sync.RWMutex
	passwordHasher  *crypto.PasswordHasher
}

// NewMemoryUserRepository 创建内存用户仓储
func NewMemoryUserRepository(userConfigs []config.UserConfig) *MemoryUserRepository {
	repo := &MemoryUserRepository{
		users:           make(map[string]*user.User),
		walletAddresses: make(map[string]*user.User),
		passwordHasher:  crypto.NewPasswordHasher(),
	}
	
	// 加载用户配置
	for _, cfg := range userConfigs {
		u := repo.createUserFromConfig(cfg)
		repo.users[u.Username] = u
		
		if u.HasWalletAddress() {
			repo.walletAddresses[strings.ToLower(u.WalletAddress)] = u
		}
	}
	
	return repo
}

// FindByUsername 根据用户名查找用户
func (r *MemoryUserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	u, ok := r.users[username]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	
	return u, nil
}

// FindByWalletAddress 根据钱包地址查找用户
func (r *MemoryUserRepository) FindByWalletAddress(ctx context.Context, address string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	address = strings.ToLower(address)
	u, ok := r.walletAddresses[address]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	
	return u, nil
}

// Save 保存用户
func (r *MemoryUserRepository) Save(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 检查用户名是否已存在
	if existing, ok := r.users[u.Username]; ok && existing.ID != u.ID {
		return user.ErrDuplicateUsername
	}
	
	// 检查钱包地址是否已存在
	if u.HasWalletAddress() {
		address := strings.ToLower(u.WalletAddress)
		if existing, ok := r.walletAddresses[address]; ok && existing.ID != u.ID {
			return user.ErrDuplicateAddress
		}
	}
	
	r.users[u.Username] = u
	
	if u.HasWalletAddress() {
		r.walletAddresses[strings.ToLower(u.WalletAddress)] = u
	}
	
	return nil
}

// Delete 删除用户
func (r *MemoryUserRepository) Delete(ctx context.Context, username string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	u, ok := r.users[username]
	if !ok {
		return user.ErrUserNotFound
	}
	
	delete(r.users, username)
	
	if u.HasWalletAddress() {
		delete(r.walletAddresses, strings.ToLower(u.WalletAddress))
	}
	
	return nil
}

// List 列出所有用户
func (r *MemoryUserRepository) List(ctx context.Context) ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	
	return users, nil
}

// createUserFromConfig 从配置创建用户
func (r *MemoryUserRepository) createUserFromConfig(cfg config.UserConfig) *user.User {
	u := user.NewUser(cfg.Username, cfg.Directory)
	
	// 设置密码
	if cfg.Password != "" {
		// 如果密码已经是加密的，直接使用
		if strings.HasPrefix(cfg.Password, "{bcrypt}") {
			u.SetPassword(cfg.Password)
		} else {
			// 否则加密密码
			hashedPassword, err := r.passwordHasher.Hash(cfg.Password)
			if err == nil {
				u.SetPassword(hashedPassword)
			}
		}
	}
	
	// 设置钱包地址
	if cfg.WalletAddress != "" {
		u.SetWalletAddress(cfg.WalletAddress)
	}
	
	// 设置权限
	if cfg.Permissions != "" {
		u.Permissions = user.ParsePermissions(cfg.Permissions)
	}
	
	// 设置规则
	for _, ruleCfg := range cfg.Rules {
		rule := &user.Rule{
			Path:        ruleCfg.Path,
			Permissions: user.ParsePermissions(ruleCfg.Permissions),
			Regex:       ruleCfg.Regex,
		}
		u.Rules = append(u.Rules, rule)
	}
	
	return u
}

