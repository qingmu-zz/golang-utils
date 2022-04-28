package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// map[string]interface{} 结构里面嵌套的可能有 map[interfacer{}]interface{}，此函数将所有的map的key都转为string类型
func Convert(my_map map[string]interface{}) (map[string]interface{}, error) {
	my_map, err := CamelMapKey(my_map)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return my_map, err
}

func CamelVecKey(my_vec []interface{}) ([]interface{}, error) {
	if !strings.HasPrefix(reflect.TypeOf(my_vec).String(), "[]") {
		return nil, errors.New("decode to vector failed")
	}
	rtn_vec := my_vec[0:0]
	for _, node := range my_vec {
		if reflect.TypeOf(node) == nil {
			continue
		}
		if strings.HasPrefix(reflect.TypeOf(node).String(), "map") {
			var res_node map[string]interface{}
			var err error
			switch v2 := node.(type) {
			case map[interface{}]interface{}:
				res_node, err = CamelMapKey(convert(v2))
			default:
				res_node, err = CamelMapKey(v2.(map[string]interface{}))
			}
			if err != nil {
				return nil, err
			}
			rtn_vec = append(rtn_vec, res_node)
		} else if strings.HasPrefix(reflect.TypeOf(node).String(), "[]") {
			res_node, err := CamelVecKey(node.([]interface{}))
			if err != nil {
				return nil, err
			}
			rtn_vec = append(rtn_vec, res_node)
		} else {
			rtn_vec = append(rtn_vec, node)
		}
	}
	return rtn_vec, nil
}

func CamelMapKey(my_map map[string]interface{}) (map[string]interface{}, error) {
	if !strings.HasPrefix(reflect.TypeOf(my_map).String(), "map") {
		return nil, errors.New("decode to map failed")
	}
	rtn_map := make(map[string]interface{})
	for k, v := range my_map {
		if reflect.TypeOf(v) == nil {
			rtn_map[fmt.Sprint(k)] = nil
			continue
		}
		if strings.HasPrefix(reflect.TypeOf(v).String(), "map") {
			var res_map map[string]interface{}
			var err error
			// res_map, err := camelMapKey(v.(map[string]interface{}))
			switch v2 := v.(type) {
			case map[interface{}]interface{}:
				res_map, err = CamelMapKey(convert(v2))
			default:
				res_map, err = CamelMapKey(v.(map[string]interface{}))
			}
			if err != nil {
				return nil, err
			}
			rtn_map[fmt.Sprint(k)] = res_map
		} else if strings.HasPrefix(reflect.TypeOf(v).String(), "[]") {
			res_node, err := CamelVecKey(v.([]interface{}))
			if err != nil {
				return nil, err
			}
			rtn_map[fmt.Sprint(k)] = res_node
		} else {
			rtn_map[fmt.Sprint(k)] = v
		}
	}
	return rtn_map, nil
}

func convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convert(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}
