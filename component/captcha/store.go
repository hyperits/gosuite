package captcha

import (
	"context"
	"time"

	"github.com/hyperits/gosuite/logger"
	"github.com/redis/go-redis/v9"
)

type CaptchaRedisStore struct {
	rc         redis.UniversalClient
	ctx        context.Context
	ExpireTime int //过期时间 秒
}

func NewCaptchaRedisStore(rc redis.UniversalClient) *CaptchaRedisStore {
	return &CaptchaRedisStore{
		rc:         rc,
		ctx:        context.Background(),
		ExpireTime: 600,
	}
}

func (rs *CaptchaRedisStore) Set(id string, digits []byte) {
	_, err := rs.rc.Set(rs.ctx, id, string(digits), time.Duration(rs.ExpireTime)*time.Second).Result()
	if err != nil {
		logger.Errorf("CaptchaRedisStore set %v failed: %v", id, err)
		return
	}
}

func (rs *CaptchaRedisStore) Get(id string, clear bool) (digits []byte) {
	v, err := rs.rc.Get(rs.ctx, id).Result()
	if err != nil {
		logger.Errorf("CaptchaRedisStore get %v failed: %v", id, err)
		return
	}
	return []byte(v)
}

func (rs *CaptchaRedisStore) Del(key string) {
	_, err := rs.rc.Del(rs.ctx, key).Result()
	if err != nil {
		logger.Errorf("CaptchaRedisStore del %v failed: %v", key, err)
	}
}
