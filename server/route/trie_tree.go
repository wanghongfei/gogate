package route

var charPosMap = make(map[rune]int)
var urlCharCount int

func init() {
	var urlCharArray = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~:/?#[]@!$&'()*+,;=%"
	for ix, char := range urlCharArray {
		charPosMap[char] = ix
	}

	urlCharCount = len(urlCharArray)
}

// 字典树
type TrieTree struct {
	root		*treeNode
}

// 创建空字典树
func NewTrieTree() *TrieTree {
	return &TrieTree{
		root: newTreeNode(0, nil),
	}
}

// 查找path对应的服务信息
// path: 请求路径
func (tree *TrieTree) Search(path string) *ServiceInfo {
	node := tree.root
	for _, char := range path {
		node = node.findSubNode(char)
		if nil == node {
			return nil
		}
	}

	return node.Data
}

// 搜索路径上遇到的第一个字符串
// path: 请求路径
func (tree *TrieTree) SearchFirst(path string) *ServiceInfo {
	node := tree.root
	for _, char := range path {
		node = node.findSubNode(char)
		if nil == node {
			return nil
		}

		if nil != node.Data {
			return node.Data
		}
	}

	return nil
}

// 添加一条path->serviceInfo映射
func (tree *TrieTree) PutString(path string, data *ServiceInfo) {
	pathRunes := []rune(path)
	LEN := len(pathRunes)

	node := tree.root
	for ix, char := range pathRunes {
		subNode := findNode(char, node.SubNodes)
		if nil == subNode {
			var newNode *treeNode
			// 是最后一个字符
			if ix == LEN - 1 {
				newNode = newTreeNode(char, data)
			} else {
				newNode = newTreeNode(char, nil)
			}

			node.addSubNode(newNode)
			node = newNode

		} else if ix == LEN - 1 {
			subNode.Data = data

		} else {
			node = subNode
		}
	}
}

func findNode(char rune, nodeList []*treeNode) *treeNode {
	if nil == nodeList {
		return nil
	}

	pos := mapPosition(char)
	return nodeList[pos]
}

type treeNode struct {
	Char		rune
	Data		*ServiceInfo

	SubNodes	[]*treeNode
}

func newTreeNode(char rune, data *ServiceInfo) *treeNode {
	node := &treeNode{
		Char: char,
		Data: data,
		SubNodes: nil,
	}

	return node
}

func mapPosition(char rune) int {
	return charPosMap[char]
}

// 添加子节点
func (node *treeNode) addSubNode(newNode *treeNode) {
	if nil == node.SubNodes {
		node.SubNodes = make([]*treeNode, urlCharCount)
	}

	position := mapPosition(newNode.Char)
	node.SubNodes[position] = newNode
}

func (node *treeNode) findSubNode(target rune) *treeNode {
	if nil == node.SubNodes {
		return nil
	}

	pos := mapPosition(target)
	return node.SubNodes[pos]
}

