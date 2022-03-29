package depstr

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"
)
type Records struct {
	Version         string            `bson:"version,omitempty"`
	Name            string            `bson:"name,omitempty"`
	DevDependencies map[string]string `bson:"devDependencies,omitempty"`
	Dependencies    map[string]string `bson:"dependencies,,omitempty"`
	Licence         string            `bson:"license,omitempty"`
}

func NewRecords() *Records {
	return &Records{
		DevDependencies: map[string]string{},
		Dependencies:    map[string]string{},
	}
}

type Metas struct {
	Version string    `bson:"version,omitempty"`
	Name    string    `bson:"name,omitempty"`
	Time    time.Time `bson:"time,omitempty"`
}



type DepTree struct {
	// 构建DepTree
	Version  Version
	Name     string
	Licence  string
	Time     time.Time
	Parent   *DepTree
	Children []*DepTree
}

/**
 * @description: 获取组件的孩子数
 * @param {*}
 * @return {*}
 */
func (dept *DepTree) DependenceDegree() int {
	return len(dept.Children)
}



/**
 * @description: 创建新的DepTree
 * @param {*DepTree} parent
 * @return {*}
 */
func NewDepTree(parent *DepTree, name, orgVersion string) *DepTree {
	dep := &DepTree{
		Name:     name,
		Version:  NewVersion(orgVersion),
		Parent:   parent,
		Children: []*DepTree{},
	}
	if parent == nil {
		return dep
	}
	parent.Children = append(parent.Children, dep)
	return dep
}


type EncodeTree struct{
	Name string `json:"name"`
	Version string `json:"version"`
	Time time.Time `json:"time"`
	License string 	`json:"licencse"`
	Index int `json:"Index"`
	ParentIndex int `json:"parentIndex"`
}


func(root *DepTree)  LevelOrder(){
	type Node struct {
		dep *DepTree
		lay int
		index int32
		parent int32
	}
	q := NewQueue()
    var index int32 =  1
	q.Push(Node{dep: root, lay: 0, index: 1, parent: 0})
	name := "%s@%sDepTree.json"
	name = fmt.Sprintf(name, root.Name, root.Version.CleanVersion)
	f, _ := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	encoder := json.NewEncoder(f)
	e := []EncodeTree{}
	for !q.Empty() {
		node := q.Pop().(Node)
		et := EncodeTree{
			Name: node.dep.Name,
			Version: node.dep.Version.CleanVersion,
			Time: node.dep.Time,
			Index: int(node.index),
			ParentIndex: int(node.parent),
			License: node.dep.Licence,
		}
		e = append(e, et)
		for _, clid := range node.dep.Children {
			q.Push(Node{
				dep: clid,
				lay: node.lay + 1,
				index: atomic.AddInt32(&index, 1),
				parent: node.index,
			})
		}
		// fmt.Printf("dep %v version %v layer %v \n", node.dep.Name, node.dep.Version.CleanVersion, node.lay)
	}
	encoder.Encode(e)
}

func(root *DepTree) PrintTree(){
	dfs(root, 0)
}

func dfs(root *DepTree, level int){
	if root == nil{
		return
	}
	s := strings.Repeat("... ", level)
	fmt.Printf("%s%v-%v>\n", s, root.Name, root.Version.CleanVersion)
	for _, clid := range root.Children{
		dfs(clid, level + 1)
	}
}