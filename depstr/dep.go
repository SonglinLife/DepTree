package depstr

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
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
 * @description: 获取外部依赖度
 * @param {*}
 * @return {*}
 */
func (dept *DepTree) DependenceDegree() int {
	return len(dept.Children)
}

/**
 * @description: 将查询mongodb的语句编码为bson
 * @param {string} sql
 * @return {*}
 */
func Sqlstr2Bson(sql string) interface{} {
	var bdoc interface{}
	bson.UnmarshalJSON([]byte(sql), &bdoc)
	return bdoc
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

func LevelOrder(root *DepTree) {
	type Node struct {
		dep *DepTree
		lay int
	}
	q := NewQueue()
	q.Push(Node{dep: root, lay: 0})
	for !q.Empty() {
		node := q.Pop().(Node)
		for _, clid := range node.dep.Children {
			q.Push(Node{
				dep: clid,
				lay: node.lay + 1,
			})
		}
		fmt.Printf("dep %v version %v layer %v \n", node.dep.Name, node.dep.Version.CleanVersion, node.lay)
	}
}
