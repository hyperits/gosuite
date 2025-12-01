package captcha

import (
	"context"
	"time"

	"github.com/hyperits/gosuite/logger"
	"github.com/redis/go-redis/v9"
)

// CaptchaRedisStore 基于 Redis 的验证码存储
// 实现 github.com/dchest/captcha.Store 接口
type CaptchaRedisStore struct {
	rc                redis.UniversalClient
	expireTimeSeconds int // 过期时间（秒）
}

// NewCaptchaRedisStore 创建基于 Redis 的验证码存储
func NewCaptchaRedisStore(rc redis.UniversalClient, expireTimeSeconds int) *CaptchaRedisStore {
	if expireTimeSeconds <= 0 {
		expireTimeSeconds = 300 // 默认 5 分钟
	}
	return &CaptchaRedisStore{
		rc:                rc,
		expireTimeSeconds: expireTimeSeconds,
	}
}

// Set 存储验证码
func (rs *CaptchaRedisStore) Set(id string, digits []byte) {
	ctx := context.Background()
	_, err := rs.rc.Set(ctx, rs.key(id), string(digits), time.Duration(rs.expireTimeSeconds)*time.Second).Result()
	if err != nil {
		logger.Errorf("CaptchaRedisStore set %v failed: %v", id, err)
	}
}

// Get 获取验证码
// clear 为 true 时，获取后立即删除验证码（防止重复使用）
func (rs *CaptchaRedisStore) Get(id string, clear bool) (digits []byte) {
	ctx := context.Background()
	v, err := rs.rc.Get(ctx, rs.key(id)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorf("CaptchaRedisStore get %v failed: %v", id, err)
		}
		return nil
	}

	// 如果 clear 为 true，获取后删除验证码
	if clear {
		rs.Del(id)
	}

	return []byte(v)
}

// Del 删除验证码
func (rs *CaptchaRedisStore) Del(key string) {
	ctx := context.Background()
	_, err := rs.rc.Del(ctx, rs.key(key)).Result()
	if err != nil {
		logger.Errorf("CaptchaRedisStore del %v failed: %v", key, err)
	}
}

// key 生成存储键，添加前缀避免键冲突
func (rs *CaptchaRedisStore) key(id string) string {
	return "captcha:" + id
}
