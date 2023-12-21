package code

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperits/gosuite/logger"
	"github.com/redis/go-redis/v9"
)

type CodeComp struct {
	rc         redis.UniversalClient
	ctx        context.Context
	ExpireTime int //过期时间 秒
	source     rand.Source
}

func NewCodeComp(rc redis.UniversalClient) *CodeComp {
	return &CodeComp{
		rc:         rc,
		ctx:        context.Background(),
		ExpireTime: 600,
		source:     rand.NewSource(time.Now().UnixNano()),
	}
}

func (c *CodeComp) GetCode(key string) string {
	min := 100000
	max := 999999
	iCode := rand.New(c.source).Intn(max-min) + min
	code := strconv.FormatInt(int64(iCode), 10)
	c.set(key, code)
	return code
}

func (c *CodeComp) set(id string, value string) {
	_, err := c.rc.Set(c.ctx, id, value, time.Duration(c.ExpireTime)*time.Second).Result()
	if err != nil {
		logger.Errorf("CodeComp set %v failed: %v", id, err)
		return
	}
}

func (c *CodeComp) Verify(key string, code string) bool {
	v, err := c.rc.Get(c.ctx, key).Result()
	if err != nil {
		logger.Errorf("CodeComp get %v failed: %v", key, err)
		return false
	}

	if v == code {
		c.Del(key)
		return true
	}

	return false
}

func (c *CodeComp) Del(key string) {
	_, err := c.rc.Del(c.ctx, key).Result()
	if err != nil {
		logger.Errorf("CodeComp del %v failed: %v", key, err)
	}
}
