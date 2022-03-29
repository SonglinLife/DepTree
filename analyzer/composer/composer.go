package composer

import (
	"context"
	"log"
	"opendep/depstr"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

/**
 * @description: 获取该组件特定版本的信息
 * @param {*depstr.DepTree} dep
 * @return {*}
 */
func Composerpkg(dep *depstr.DepTree) bool {

	// enrich信息
	records, fl := ComposerRecords(dep.Name, dep.Version.CleanVersion)

	if !fl {
		return fl
	}
	dep.Licence = strings.Join(records.Licence, ",")

	// 添加到子依赖
	// for name, orgVersion := range records.DevDependencies {
	// 	depstr.NewDepTree(dep, name, orgVersion)
	// }
	for name, orgVersion := range records.Require {
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
		constraints := []*semver.Constraints{}
		for _, constraint := range strings.Split(cdep.Version.Org, "|") {
			constraint = strings.TrimSpace(constraint)
			if !strings.EqualFold(constraint, "*") {
				constraint = strings.ReplaceAll(constraint, "*", "x")

			}
			if c, err := semver.NewConstraint(constraint); err == nil {
				constraints = append(constraints, c)
			}
		}

		metas := ComposerMetas(cdep.Name)

		var max *semver.Version
		var v = cdep.Version.Org
		var t = cdep.Time
		for _, meta := range metas {
			if strings.EqualFold(meta.Name, "php") {
				continue
			}
			if version, err := semver.NewVersion(meta.Version); err == nil {
				for _, constraint := range constraints {
					if constraint.Check(version) {
						if max == nil || max.LessThan(version) {
							max = version
							v = meta.Name
							t = meta.Time
						}
					}
				}
			}

		}

		cdep.Version.CleanVersion = v
		cdep.Time = t

	}
}

func BuildDepTree(root *depstr.DepTree) bool {
	ch := make(chan bool, 1)
	ctx, c := context.WithTimeout(context.Background(), 15*time.Second)
	// 限制建立依赖树的时间不超过 15 s
	defer c()
	go func() {
		ch <- buildTree(root)
	}()
	select {
	case <-ctx.Done():
		log.Println("out of time")
		return false
	case <-ch:
		return true
	}
}

func buildTree(root *depstr.DepTree) bool {
	q := depstr.NewQueue()
	q.Push(root)
	for !q.Empty() {
		node := q.Pop().(*depstr.DepTree)
		Composerpkg(node)  // 获取该节点的子依赖信息
		CleanClidren(node) // 解析为最大满足的版本号
		for _, clid := range node.Children {
			q.Push(clid)
		}
	}
	return true
}
