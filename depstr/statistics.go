package depstr

import "time"

func (root *DepTree) External() float32 {
	if root == nil {
		return 0.0
	}
	if len(root.Children) == 0 {
		return 0.0
	}
	n := float32(len(root.Children))
	var external = n
	for _, child := range root.Children {
		external += child.External() / n
	}
	return external
}

func (root *DepTree) MTTU() float32 {
	var sum float32 = 0.0
	var cnt = 0.0
	for _, clid := range root.Children {
		
		if clid.Version.CleanVersion != "" { // 没有查到
			sum += float32(root.Time.Sub(clid.Time) / (time.Hour * 24))
			cnt++
		}
	}
	return sum / float32(cnt) 
}