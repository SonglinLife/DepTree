package mttu

import (
	"opendep/analyzer/npm"
	"opendep/depstr"
	"time"
	"github.com/Masterminds/semver"
)

/**
 * @description:
 * @param {*depstr.DepTree} root
 * @return {*}
 */
func BuildMttuTree(root *depstr.DepTree) bool {
	ch := make(chan bool, 1)
	s := make(chan bool, 1)

	// 限制建立依赖树的时间不超过15s
	go func() {
		select {
		case ch <- buildTree(root):
		case <-s:
		}
	}()

	for {
		select {
		case <-time.After(15 * time.Second):
			s <- true
			return false
		case <-ch:
			return true
		}
	}
}

/**
 * @description: mttu 树的建立比较特殊，不需要
 * @param {*depstr.DepTree} root
 * @return {*}
 */
func buildTree(root *depstr.DepTree) bool {
	npm.Npmpkg(root) // 获取该节点的子依赖信息
	for _, clid := range root.Children {
		clid.Time = root.Time
		clid.Version.CleanVersion = ""
	}
	for _, clid := range root.Children {
		if constraint, err := semver.NewConstraint(clid.Version.Org); err == nil {
			cv, ti := GetMaxVersion(clid.Name, root.Time, constraint)
			if cv != "" {
				clid.Version.CleanVersion = cv
				clid.Time = ti
			} else {
				clid.Time = root.Time
				clid.Version.CleanVersion = cv
			}
		}
	}
	return true
}

/**
 * @description: 查询限制时间范围内最大的mmtu。
 * @param {string} depName
 * @param {time.Time} t
 * @param {*semver.Constraints} constraint
 * @return {*}
 */
func GetMaxVersion(depName string, t time.Time, constraint *semver.Constraints) (string, time.Time) {
	v := npm.NpmMetas(depName)
	var max *semver.Version
	var ti time.Time
	var ver string = ""
	for _, depVersion := range v {
		if version, err := semver.NewVersion(depVersion.Version); err == nil && depVersion.Time.Before(t) {
			if constraint.Check(version) {
				if max == nil || max.LessThan(version) {
					max = version
					ti = depVersion.Time
					ver = depVersion.Version
				}
			}
		}
	}
	return ver, ti
}
