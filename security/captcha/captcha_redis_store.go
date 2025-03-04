package captcha

import (
	"context"
	"time"

	"github.com/hyperits/gosuite/log"
	"github.com/redis/go-redis/v9"
)

type CaptchaRedisStore struct {
	rc                redis.UniversalClient
	ctx               context.Context
	expireTimeSeconds int //过期时间 秒
}

func NewCaptchaRedisStore(rc redis.UniversalClient, expireTimeSeconds int) *CaptchaRedisStore {
	return &CaptchaRedisStore{
		rc:                rc,
		ctx:               context.Background(),
		expireTimeSeconds: expireTimeSeconds,
	}
}

func (rs *CaptchaRedisStore) Set(id string, digits []byte) {
	_, err := rs.rc.Set(rs.ctx, id, string(digits), time.Duration(rs.expireTimeSeconds)*time.Second).Result()
	if err != nil {
		log.Errorf("CaptchaRedisStore set %v failed: %v", id, err)
		return
	}
}

func (rs *CaptchaRedisStore) Get(id string, clear bool) (digits []byte) {
	v, err := rs.rc.Get(rs.ctx, id).Result()
	if err != nil {
		log.Errorf("CaptchaRedisStore get %v failed: %v", id, err)
		return
	}
	return []byte(v)
}

func (rs *CaptchaRedisStore) Del(key string) {
	_, err := rs.rc.Del(rs.ctx, key).Result()
	if err != nil {
		log.Errorf("CaptchaRedisStore del %v failed: %v", key, err)
	}
}
