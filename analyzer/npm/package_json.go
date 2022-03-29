package npm

import (
	"context"
	"log"
	"opendep/depstr"
	"sort"
	"time"

	"github.com/Masterminds/semver/v3"
)

/**
 * @description: 获取该组件特定版本的信息
 * @param {*depstr.DepTree} dep
 * @return {*}
 */
func Npmpkg(dep *depstr.DepTree) bool {

	// enrich信息
	records, fl := NpmRecords(dep.Name, dep.Version.CleanVersion)

	if !fl {
		return fl
	}
	dep.Licence = records.Licence

	// 添加到子依赖
	// for name, orgVersion := range records.DevDependencies {
	// 	depstr.NewDepTree(dep, name, orgVersion)
	// }
	for name, orgVersion := range records.Dependencies {
		depstr.NewDepTree(dep, name, orgVersion)
	}
	return fl
}

/**
 * @description: 将package.json 中dep的semver版本号转换为最大满足要求的版本号
 * @param {*depstr.DepTree} dep
 * @return {*}
 */
func CleanClidren(dep *depstr.DepTree) {

	for _, cdep := range dep.Children {

		constraint, err := semver.NewConstraint(cdep.Version.Org)
		if err != nil {
			return
		}
		metas := NpmMetas(cdep.Name)
		sort.Slice(metas, func(i, j int) bool {
			return metas[i].Time.After(metas[j].Time)
		})
		var max *semver.Version
		var ver = cdep.Version.Org
		var t time.Time

		for _, meta := range metas {

			version, err := semver.NewVersion(meta.Version)
			if err != nil {
				continue
			}
			// fmt.Println(version.Original())
			if constraint.Check(version) {
				if max == nil || max.LessThan(version) {
					max = version
					ver = meta.Version
					t = meta.Time
				}

			}
		}

		cdep.Version.CleanVersion = ver
		cdep.Time = t
	}
}


func BuildDepTree(root *depstr.DepTree) bool {
	ch := make(chan bool, 1)
	s  := make(chan bool, 1)
	ctx, c := context.WithTimeout(context.Background(), 30*time.Second)
	// 限制建立依赖树的时间不超过5s
	defer c()
	go func() {
		select{
		case ch <- buildTree(root):
		case <- s:
		}
		
	}()
	for{
		select {
		case <-ctx.Done():
			log.Printf("name %v version %v out of time", root.Name, root.Version)
			// time.Sleep(10 *time.Second)
			// ctx, c = context.WithTimeout(context.Background(), 10*time.Second)
			// 限制建立依赖树的时间不超过5s

			s <- true
			return false
		
		case <-ch:
			return true
		}
	}


}

func buildTree(root *depstr.DepTree) bool {
	q := depstr.NewQueue()

	type Node struct {
		deptree *depstr.DepTree
		layer int
	}
	q.Push(Node{
		deptree: root,
		layer: 1,
	})
	noLoop := map[string]int{}
	noLoop[root.Name] = 1
	for !q.Empty() {
		node := q.Pop().(Node)
		Npmpkg(node.deptree)       // 获取该节点的子依赖信息
		CleanClidren(node.deptree) // 解析为最大满足的版本号
		for _, clid := range node.deptree.Children {
			if noLoop[clid.Name] ==0 || noLoop[clid.Name] == node.layer + 1{
				// 不准出现依赖回溯
				q.Push(Node{
					deptree: clid,
					layer: node.layer + 1,
				})
				noLoop[clid.Name] = node.layer + 1
			}
				
		}
	}
	return true
}

