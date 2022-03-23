package npm

import (
	"opendep/depstr"
	"testing"
)

func TestBuild(t *testing.T) {
	name := "field-descriptions"
	version := "1.0.7"

	root := depstr.NewDepTree(nil, name, version)

	BuildDepTree(root)
	
	// depstr.LevelOrder(root)
}