package npm

import (
	"fmt"
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
		fmt.Println(cdep.Version.Org)
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
			if err != nil{
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

func BuildDepTree(root *depstr.DepTree) {

	q := depstr.NewQueue()
	q.Push(root)
	for !q.Empty() {
		node := q.Pop().(*depstr.DepTree)
		Npmpkg(node)       // 获取该节点的子依赖信息
		CleanClidren(node) // 解析为最大满足的版本号
		for _, clid := range node.Children {
			q.Push(clid)
		}
	}

}
