package server

import (
	"github.com/alecthomas/log4go"
)

// 注册过滤器, 追加到末尾
func (serv *Server) AppendPreFilter(pre *PreFilter) {
	if serv.isStarted {
		log4go.Warn("cannot change filters after server started")
	}

	log4go.Info("append pre filter: %s", pre.Name)
	serv.preFilters = append(serv.preFilters, pre)
}

// 注册过滤器, 追加到末尾
func (serv *Server) AppendPostFilter(post *PostFilter) {
	if serv.isStarted {
		log4go.Warn("cannot change filters after server started")
	}

	log4go.Info("append post filter: %s", post.Name)
	serv.postFilters = append(serv.postFilters, post)
}

func (serv *Server) ExportAllPreFilters() []*PreFilter {
	result := make([]*PreFilter, len(serv.preFilters))
	copy(result, serv.preFilters)

	return result
}

func (serv *Server) ExportAllPostFilters() []*PostFilter {
	result := make([]*PostFilter, len(serv.postFilters))
	copy(result, serv.postFilters)

	return result
}

// 在指定前置过滤器的后面添加
func (serv *Server) InsertPreFilterBehind(filterName string, filter *PreFilter) bool {
	if serv.isStarted {
		log4go.Warn("cannot change filters after server started")
	}

	log4go.Info("insert pre filter: %s", filter.Name)

	targetIdx := serv.getPreFilterIndex(filterName)
	if -1 == targetIdx {
		return false
	}

	rearIdx := targetIdx + 1
	rear := append([]*PreFilter{}, serv.preFilters[rearIdx:]...)
	serv.preFilters = append(serv.preFilters[0:rearIdx], filter)
	serv.preFilters = append(serv.preFilters, rear...)


	return true
}

// 在指定后置过滤器的后面添加
// filterName: 在此过滤器后面添加filter, 如果要在队头添加, 则使用空字符串
// filter: 过滤器对象
func (serv *Server) InsertPostFilterBehind(filterName string, filter *PostFilter) bool {
	if serv.isStarted {
		log4go.Warn("cannot change filters after server started")
	}

	log4go.Info("insert post filter: %s", filter.Name)

	targetIdx := serv.getPostFilterIndex(filterName)
	if -1 == targetIdx {
		return false
	}

	rearIdx := targetIdx + 1
	rear := append([]*PostFilter{}, serv.postFilters[rearIdx:]...)
	serv.postFilters = append(serv.postFilters[0:rearIdx], filter)
	serv.postFilters = append(serv.postFilters, rear...)

	return true
}

// 在最头部添加前置过滤器
func (serv *Server) InsertPreFilterAhead(filter *PreFilter) {
	if serv.isStarted {
		log4go.Warn("cannot change filters after server started")
	}

	log4go.Info("insert pre filter: %s", filter.Name)

	newFilterSlice := make([]*PreFilter, 0, 1 + len(serv.preFilters))
	newFilterSlice = append(newFilterSlice, filter)
	newFilterSlice = append(newFilterSlice, serv.preFilters...)

	serv.preFilters = newFilterSlice
}

// 在最头部添加后置过滤器
func (serv *Server) InsertPostFilterAhead(filter *PostFilter) {
	if serv.isStarted {
		log4go.Warn("cannot change filters after server started")
	}

	log4go.Info("insert post filter: %s", filter.Name)

	newFilterSlice := make([]*PostFilter, 0, 1 + len(serv.postFilters))
	newFilterSlice = append(newFilterSlice, filter)
	newFilterSlice = append(newFilterSlice, serv.postFilters...)

	serv.postFilters = newFilterSlice
}

func (serv *Server) ensurePreFilterCap(neededSpace int) {
	currentCap := cap(serv.preFilters)
	currentLen := len(serv.preFilters)
	leftSpace := currentCap - currentLen

	if leftSpace < neededSpace {
		newCap := currentCap + (neededSpace - leftSpace) + 3

		oldFilters := serv.preFilters
		serv.preFilters = make([]*PreFilter, 0, newCap)
		copy(serv.preFilters, oldFilters)
	}
}

func (serv *Server) getPreFilterIndex(name string) int {
	if nil == serv.preFilters {
		return -1
	}

	for ix, f := range serv.preFilters {
		if f.Name == name {
			return ix
		}
	}

	return -1
}

func (serv *Server) getPostFilterIndex(name string) int {
	if nil == serv.preFilters {
		return -1
	}

	for ix, f := range serv.postFilters {
		if f.Name == name {
			return ix
		}
	}

	return -1
}

