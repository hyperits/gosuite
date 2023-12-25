package converter

import (
	"encoding/json"
)

// ObjectToMap 将一个任意类型的对象转换为map[string]interface{}类型
// 参数obj是待转换的对象
// 返回值是一个指向map[string]interface{}的指针和一个error
func ObjectToMap(obj interface{}) (*map[string]interface{}, error) {

	bits, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]interface{})
	err = json.Unmarshal(bits, &mp)
	if err != nil {
		return nil, err
	}

	return &mp, nil
}
