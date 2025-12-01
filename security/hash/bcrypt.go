package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost 是 bcrypt 的默认成本因子
// 值为 10，在安全性和性能之间取得平衡
// 可通过 BcryptHashPasswordWithCost 自定义成本
const DefaultCost = bcrypt.DefaultCost

// BcryptHashPassword 使用 bcrypt 算法对密码进行哈希处理
// 使用默认成本因子（10），适合大多数场景
func BcryptHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// BcryptHashPasswordWithCost 使用指定的成本因子对密码进行哈希处理
// cost 值范围：4-31，推荐值：10-14
// 更高的成本会增加计算时间，提高安全性，但也会增加 CPU 负担
func BcryptHashPasswordWithCost(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// BcryptMatchPassword 验证密码是否与哈希值匹配
func BcryptMatchPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
