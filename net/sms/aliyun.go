package sms

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type AliyunSmsConfig struct {
	Region       string `yaml:"region"`
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
}

type AliSmsResp struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

type AliyunSmsClient struct {
	conf *AliyunSmsConfig
}

func NewAliyunSmsClient(conf *AliyunSmsConfig) *AliyunSmsClient {
	return &AliyunSmsClient{
		conf: conf,
	}
}

func (c *AliyunSmsClient) SendCode(mobile string, code string) error {

	conf := c.conf
	client, err := sdk.NewClientWithAccessKey(conf.Region, conf.AccessKey, conf.SecretKey)
	if err != nil {
		return err
	}

	request := requests.NewCommonRequest()                           // 构造一个公共请求
	request.Method = "POST"                                          // 设置请求方式
	request.Product = "Ecs"                                          // 指定产品
	request.Scheme = "https"                                         // https | http
	request.Domain = "dysmsapi.aliyuncs.com"                         // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2017-05-25"                                   // 指定产品版本
	request.ApiName = "SendSms"                                      // 指定接口名
	request.QueryParams["RegionId"] = conf.Region                    // 地区
	request.QueryParams["PhoneNumbers"] = mobile                     // 手机号
	request.QueryParams["SignName"] = conf.SignName                  // 阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = conf.TemplateCode          // 阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + code + "}" // 短信模板中的验证码内容 自己生成

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	var resp AliSmsResp
	if err := json.Unmarshal(response.GetHttpContentBytes(), &resp); err != nil {
		return err
	}

	if resp.Message != "OK" {
		return fmt.Errorf("aliyun resp not OK")
	}
	return nil
}

func (c *AliyunSmsClient) SendRaw(mobile string, templateCode string, templateParam string) error {

	conf := c.conf
	client, err := sdk.NewClientWithAccessKey(conf.Region, conf.AccessKey, conf.SecretKey)
	if err != nil {
		return err
	}

	request := requests.NewCommonRequest()               // 构造一个公共请求
	request.Method = "POST"                              // 设置请求方式
	request.Product = "Ecs"                              // 指定产品
	request.Scheme = "https"                             // https | http
	request.Domain = "dysmsapi.aliyuncs.com"             // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2017-05-25"                       // 指定产品版本
	request.ApiName = "SendSms"                          // 指定接口名
	request.QueryParams["RegionId"] = conf.Region        // 地区
	request.QueryParams["PhoneNumbers"] = mobile         // 手机号
	request.QueryParams["SignName"] = conf.SignName      // 阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = templateCode   // 阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = templateParam // 短信模板中的验证码内容 自己生成

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	var resp AliSmsResp
	if err := json.Unmarshal(response.GetHttpContentBytes(), &resp); err != nil {
		return err
	}

	if resp.Message != "OK" {
		return fmt.Errorf("aliyun resp not OK")
	}
	return nil
}
