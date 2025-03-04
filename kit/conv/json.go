package conv

import (
	"encoding/json"

	"github.com/hyperits/gosuite/kit/debug"
	"github.com/hyperits/gosuite/log"
)

// ObjectToJsonString 将任意类型的v转换为格式化的JSON字符串
// 参数v为任意类型
// 返回值为格式化的JSON字符串
func ObjectToJsonString(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.ErrorRTf(debug.GetCurrentFunctionInfo(), "ToJsonString error: %v", err)
		return ""
	}

	return string(bytes)
}

// ObjectToJsonStringErr 将一个任意类型的v转换成json格式的字符串
// 参数v是任意类型
// 返回值第一个是转换后的json字符串，第二个是转换过程中出现的错误，如果没有错误则为nil
func ObjectToJsonStringErr(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
