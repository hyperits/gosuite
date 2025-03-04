package verify

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperits/gosuite/log"
	"github.com/redis/go-redis/v9"
)

type CodeClient struct {
	client            redis.UniversalClient
	ctx               context.Context
	expireTimeSeconds int // 过期时间（秒）
	source            rand.Source
}

func NewCodeClient(client redis.UniversalClient, expireTimeSeconds int) *CodeClient {
	return &CodeClient{
		client:            client,
		ctx:               context.Background(),
		expireTimeSeconds: expireTimeSeconds,
		source:            rand.NewSource(time.Now().UnixNano()),
	}
}

func (c *CodeClient) GenerateCode(key string) string {
	min := 100000
	max := 999999
	randomCode := rand.New(c.source).Intn(max-min) + min
	code := strconv.FormatInt(int64(randomCode), 10)
	c.setCode(key, code)
	return code
}

func (c *CodeClient) setCode(key string, value string) {
	_, err := c.client.Set(c.ctx, key, value, time.Duration(c.expireTimeSeconds)*time.Second).Result()
	if err != nil {
		log.Errorf("CodeClient setCode %v failed: %v", key, err)
		return
	}
}

func (c *CodeClient) VerifyCode(key string, code string) bool {
	storedCode, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		log.Errorf("CodeClient verifyCode %v failed: %v", key, err)
		return false
	}

	if storedCode == code {
		c.deleteCode(key)
		return true
	}

	return false
}

func (c *CodeClient) deleteCode(key string) {
	_, err := c.client.Del(c.ctx, key).Result()
	if err != nil {
		log.Errorf("CodeClient deleteCode %v failed: %v", key, err)
	}
}
