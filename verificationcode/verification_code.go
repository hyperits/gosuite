package verificationcode

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperits/gosuite/logger"
	"github.com/redis/go-redis/v9"
)

type VerificationCodeComponent struct {
	client     redis.UniversalClient
	ctx        context.Context
	expireTime int // 过期时间（秒）
	source     rand.Source
}

func NewVerificationCodeComponent(client redis.UniversalClient) *VerificationCodeComponent {
	return &VerificationCodeComponent{
		client:     client,
		ctx:        context.Background(),
		expireTime: 600,
		source:     rand.NewSource(time.Now().UnixNano()),
	}
}

func (c *VerificationCodeComponent) GenerateCode(key string) string {
	min := 100000
	max := 999999
	randomCode := rand.New(c.source).Intn(max-min) + min
	code := strconv.FormatInt(int64(randomCode), 10)
	c.setCode(key, code)
	return code
}

func (c *VerificationCodeComponent) setCode(key string, value string) {
	_, err := c.client.Set(c.ctx, key, value, time.Duration(c.expireTime)*time.Second).Result()
	if err != nil {
		logger.Errorf("VerificationCodeComponent setCode %v failed: %v", key, err)
		return
	}
}

func (c *VerificationCodeComponent) VerifyCode(key string, code string) bool {
	storedCode, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		logger.Errorf("VerificationCodeComponent verifyCode %v failed: %v", key, err)
		return false
	}

	if storedCode == code {
		c.deleteCode(key)
		return true
	}

	return false
}

func (c *VerificationCodeComponent) deleteCode(key string) {
	_, err := c.client.Del(c.ctx, key).Result()
	if err != nil {
		logger.Errorf("VerificationCodeComponent deleteCode %v failed: %v", key, err)
	}
}
