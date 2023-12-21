package captcha

import (
	"bytes"
	"encoding/base64"

	"github.com/dchest/captcha"
	"github.com/hyperits/gosuite/logger"
)

const (
	DefaultLen    = 4   // 验证码默认长度 4 位数字
	DefaultWidth  = 200 // 验证码默认宽度
	DefaultHeight = 100 // 验证码默认高度
)

type CaptchaComp struct {
	CaptchaLen    int `json:"captcha_len"`    //验证码长度
	CaptchaWidth  int `json:"captcha_width"`  //验证码图片宽度
	CaptchaHeight int `json:"captcha_height"` //验证码图片高度
	store         *CaptchaRedisStore
}

func NewCaptchaComp(store *CaptchaRedisStore, len int, width int, height int) *CaptchaComp {
	comp := &CaptchaComp{
		CaptchaLen:    len,
		CaptchaWidth:  width,
		CaptchaHeight: height,
		store:         store,
	}

	captcha.SetCustomStore(store)
	comp.verifyParma()

	return comp
}

// GetCaptcha 获取一个验证码 imageData 存放base64之后的图片信息
func (c *CaptchaComp) GetCaptcha() (captchaId string, imageData string, err error) {
	captchaId = captcha.NewLen(c.CaptchaLen)
	var image bytes.Buffer
	err = captcha.WriteImage(&image, captchaId, c.CaptchaWidth, c.CaptchaHeight)
	imageData = base64.StdEncoding.EncodeToString(image.Bytes())
	if err != nil {
		return
	}
	return
}

// VerifyCaptcha 验证是否正确 digits 前端传过来的数字字符串验证码
// 验证成功删除redis_key
func (c *CaptchaComp) VerifyCaptcha(captchaId string, digits string) bool {
	res := captcha.VerifyString(captchaId, digits)
	if !res {
		return false
	}
	go func() {
		c.store.Del(captchaId)
	}()
	return true
}

// 验证参数
func (c *CaptchaComp) verifyParma() {
	if c.CaptchaLen <= 0 {
		c.CaptchaLen = DefaultLen
		logger.Warnf("Invalid captcha len, use default [%v]", DefaultLen)
	}
	if c.CaptchaWidth <= 0 {
		c.CaptchaWidth = DefaultWidth
		logger.Warnf("Invalid captcha width, use default [%v]", DefaultWidth)
	}
	if c.CaptchaHeight <= 0 {
		c.CaptchaHeight = DefaultHeight
		logger.Warnf("Invalid captcha height, use default [%v]", DefaultHeight)
	}
}
