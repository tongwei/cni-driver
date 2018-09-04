package utils

import (
	"strings"

	"github.com/rancher/go-rancher-metadata/metadata"
	"reflect"
)

const (
	hostLabelKeyword = "__host_label__"
)

func toSlice(arg interface{}) (rst []interface{}, ok bool) {
	s := reflect.ValueOf(arg)
	if s.Kind() != reflect.Slice {
		return
	}
	l := s.Len()
	rst = make([]interface{}, l)
	for i := 0; i < l; i++ {
		rst[i] = s.Index(i).Interface()
	}
	ok = true
	return
}

// UpdateCNIConfigByKeywords takes in the given CNI config, replaces the rancher
// specific keywords with the appropriate values.
func UpdateCNIConfigByKeywords(config interface{}, host metadata.Host) interface{} {
	props, isMap := config.(map[string]interface{})
	if !isMap {
		sliceProps, isSlice := toSlice(config)
		if !isSlice {
			return config
		}
		for idx, v := range sliceProps {
			sliceProps[idx] = UpdateCNIConfigByKeywords(v, host)
		}
		return sliceProps
	}

	for aKey, aValue := range props {
		if v, isString := aValue.(string); isString {
			if strings.HasPrefix(v, hostLabelKeyword) {
				props[aKey] = ""
				splits := strings.SplitN(v, ":", 2)
				if len(splits) > 1 {
					label := strings.TrimSpace(splits[1])
					labelValue := host.Labels[label]
					if labelValue != "" {
						props[aKey] = labelValue
					}
				}
			}
		} else {
			props[aKey] = UpdateCNIConfigByKeywords(aValue, host)
		}
	}

	return props
}
