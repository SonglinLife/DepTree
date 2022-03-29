package npm

import (
	"fmt"
	"testing"
	"time"
)

func TestSql(t *testing.T) {
	ti := VersionTime("field-descriptions", "1.0.7")
	if ti!=nil{
		t := time.Now()
		fmt.Println(t)
	}
}
