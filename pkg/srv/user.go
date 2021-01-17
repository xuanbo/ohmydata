package srv

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/xuanbo/ohmydata/pkg/cache"
	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/entity"

	"gorm.io/gorm"
)

const (
	cacheTTL = 5 * time.Minute
)

// User 用户服务
type User struct {
	db *gorm.DB
}

// NewUser 创建实例
func NewUser() *User {
	return &User{db: db.DB}
}

// Username 查询用户
func (u *User) Username(ctx context.Context, username string) (*entity.User, error) {
	var (
		user entity.User
		key  = "ohmydata:user:username:" + username
		err  error
	)
	if err = cache.Get(ctx, key, &user); errors.Is(err, redis.Nil) {
		// 查询db
		if err = u.db.WithContext(ctx).Where("username = ?", username).Find(&user).Error; err != nil {
			return nil, err
		}
		if user.ID == "" {
			return nil, nil
		}
		// 写入缓存
		cache.Set(ctx, key, &user, cacheTTL)
	}
	return &user, err
}

// Login 登录
func (u *User) Login(ctx context.Context, user *entity.User) (string, error) {
	if user.Username == "" {
		return "", errors.New("用户名不能为空")
	}
	if user.Password == "" {
		return "", errors.New("密码不能为空")
	}
	s, err := u.Username(ctx, user.Username)
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", errors.New("用户不存在")
	}
	if s.Password != user.Password {
		return "", errors.New("密码错误")
	}
	// Create the Claims
	claims := &jwt.StandardClaims{
		// 30分钟过期
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
	}
	token := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": jwt.SigningMethodHS256.Alg(),
			// 用户信息
			"userId":   s.ID,
			"userName": s.Name,
			"username": s.Username,
		},
		Claims: claims,
		Method: jwt.SigningMethodHS256,
	}
	return token.SignedString([]byte("secret"))
}
