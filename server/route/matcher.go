package route

import "strings"

type PathMatcher struct {
	routeMap		map[string]*ServiceInfo
	routeTrieTree	*TrieTree
}

func (matcher *PathMatcher) Match(path string) *ServiceInfo {
	// 如果大于3个token则使用TrieTree匹配提高性能
	if strings.Count(path, "/") >= 3 {
		return matcher.matchByTree(path)
	}

	// 使用切token的方式匹配
	return matcher.matchByToken(path)
}

func (matcher *PathMatcher) matchByTree(path string) *ServiceInfo {
	return matcher.routeTrieTree.SearchFirst(path)
}

func (matcher *PathMatcher) matchByToken(path string) *ServiceInfo {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if "/" == path {
		path = "//"
	}

	// 以/为分隔符, 从后向前匹配
	// 每次循环都去掉最后一个/XXXX节点
	term := path
	for {
		lastSlash := strings.LastIndex(term, "/")
		if -1 == lastSlash {
			break
		}

		matchTerm := term[0:lastSlash]
		term = matchTerm

		if "" == matchTerm {
			matchTerm = "/"
		}

		appId, exist := matcher.routeMap[matchTerm]
		if exist {
			return appId
		}
	}

	return nil

}
