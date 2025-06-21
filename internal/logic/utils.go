package logic

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/maxliu9403/common/logger"
	"gorm.io/gorm"
)

// GenerateSecureToken returns a secure random string of given byte length.
// 例如传入 32，返回 base64 后长度约为 44 个字符
func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

// RetryTransaction 事物重试
func RetryTransaction(db *gorm.DB, fn func(tx *gorm.DB) error, maxRetry int) error {
	var err error
	for i := 0; i < maxRetry; i++ {
		err = db.Transaction(fn)
		if err == nil {
			return nil
		}
		logger.Warnf("事务第 %d 次重试失败: %s", i+1, err.Error())
		time.Sleep(time.Duration(100*i) * time.Millisecond)
	}
	return fmt.Errorf("事务重试失败: %w", err)
}
