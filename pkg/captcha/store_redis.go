package captcha

import (
	"errors"
	"gohub/pkg/redis"

	"gohub/pkg/app"
	"gohub/pkg/config"
	"time"
)

// RedisStore实现base64Captcha.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	keyPrefix   string
}

// Set实现base64Captcha.Store interface的Set方法
func (s *RedisStore) Set(key string, value string) error {
	ExpireTime := time.Minute * time.Duration(config.GetInt64("captcha.expire_time"))
	// 方便本地开发调试
	if app.IsLocal() {
		ExpireTime = time.Minute * time.Duration(config.GetInt64("captcha.debug_expire_time"))
	}
	if ok := s.RedisClient.Set(s.keyPrefix+key, value, ExpireTime); !ok {
		return errors.New("无法存储图片验证码答案")
	}
	return nil
}

// Get 实现base64Captcha.Store interface的Get方法
func (s *RedisStore) Get(key string, clear bool) string {
	key = s.keyPrefix + key
	val := s.RedisClient.Get(key)
	if clear {
		s.RedisClient.Del(key)
	}
	return val
}

// Verify 实现base64Captcha.Store interface 的Verify方法
func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
