package route

import (
	"fmt"
	"testing"
)

func TestTrieTree_Search(t *testing.T) {
	tree := NewTrieTree()


	data := &ServiceInfo{}
	data.Id = "serviceA"
	tree.PutString("/", data)
	tree.PutString("abcde", data)
	fmt.Println(tree.SearchFirst("/abc"))

	data = tree.Search("abcde")
	fmt.Println(data.Id)

	data = tree.SearchFirst("abcdefgasdf")
	fmt.Println(data.Id)

}
