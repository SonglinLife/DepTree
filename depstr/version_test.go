package depstr

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestNewversion(t *testing.T) {
	v, _ := semver.NewVersion("1.2.3")
	fmt.Println(v)
}

