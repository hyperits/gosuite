package verify

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/hyperits/gosuite/log"
	"github.com/redis/go-redis/v9"
)

// CodeClient 验证码客户端，用于生成和验证验证码
type CodeClient struct {
	client            redis.UniversalClient
	expireTimeSeconds int // 过期时间（秒）
	codeLength        int // 验证码长度，默认 6 位
}

// NewCodeClient 创建验证码客户端
// client: Redis 客户端
// expireTimeSeconds: 验证码过期时间（秒）
func NewCodeClient(client redis.UniversalClient, expireTimeSeconds int) *CodeClient {
	return &CodeClient{
		client:            client,
		expireTimeSeconds: expireTimeSeconds,
		codeLength:        6,
	}
}

// NewCodeClientWithLength 创建指定长度的验证码客户端
func NewCodeClientWithLength(client redis.UniversalClient, expireTimeSeconds int, codeLength int) *CodeClient {
	if codeLength < 4 {
		codeLength = 4
	}
	if codeLength > 10 {
		codeLength = 10
	}
	return &CodeClient{
		client:            client,
		expireTimeSeconds: expireTimeSeconds,
		codeLength:        codeLength,
	}
}

// GenerateCode 生成验证码并存储到 Redis
// 使用 crypto/rand 生成安全的随机数
func (c *CodeClient) GenerateCode(key string) (string, error) {
	return c.GenerateCodeWithContext(context.Background(), key)
}

// GenerateCodeWithContext 使用指定 context 生成验证码
func (c *CodeClient) GenerateCodeWithContext(ctx context.Context, key string) (string, error) {
	code, err := c.generateSecureCode()
	if err != nil {
		return "", fmt.Errorf("generate secure code failed: %w", err)
	}

	if err := c.setCode(ctx, key, code); err != nil {
		return "", err
	}
	return code, nil
}

// generateSecureCode 使用 crypto/rand 生成安全的随机验证码
func (c *CodeClient) generateSecureCode() (string, error) {
	// 计算范围：例如 6 位数字，范围是 100000 到 999999
	min := int64(1)
	for i := 1; i < c.codeLength; i++ {
		min *= 10
	}
	max := min * 10

	// 使用 crypto/rand 生成安全随机数
	n, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return "", err
	}

	code := n.Int64() + min
	return fmt.Sprintf("%0*d", c.codeLength, code), nil
}

func (c *CodeClient) setCode(ctx context.Context, key string, value string) error {
	_, err := c.client.Set(ctx, key, value, time.Duration(c.expireTimeSeconds)*time.Second).Result()
	if err != nil {
		log.Errorf("CodeClient setCode %v failed: %v", key, err)
		return fmt.Errorf("set code failed: %w", err)
	}
	return nil
}

// VerifyCode 验证验证码是否正确
// 验证成功后会自动删除验证码
func (c *CodeClient) VerifyCode(key string, code string) bool {
	return c.VerifyCodeWithContext(context.Background(), key, code)
}

// VerifyCodeWithContext 使用指定 context 验证验证码
func (c *CodeClient) VerifyCodeWithContext(ctx context.Context, key string, code string) bool {
	storedCode, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("CodeClient verifyCode %v failed: %v", key, err)
		}
		return false
	}

	if storedCode == code {
		c.deleteCode(ctx, key)
		return true
	}

	return false
}

func (c *CodeClient) deleteCode(ctx context.Context, key string) {
	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		log.Errorf("CodeClient deleteCode %v failed: %v", key, err)
	}
}
