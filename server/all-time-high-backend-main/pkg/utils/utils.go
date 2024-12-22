package utils

import (
	"strconv"
)

// ParsePage 解析页码字符串，转换为 int。若转换失败则返回错误
func ParsePage(pageStr string) (int, error) {
	if pageStr == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1, err
	}
	return page, nil
}

// ParsePageSize 解析每页大小字符串，转换为 int。若转换失败则返回错误
func ParsePageSize(pageSizeStr string) (int, error) {
	if pageSizeStr == "" {
		return 4, nil
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return 4, err
	}
	return pageSize, nil
}
