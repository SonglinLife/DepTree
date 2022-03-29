package sample

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestSample(t *testing.T) {
	SampleMetas("composer_metas")
}

func TestParse(t *testing.T) {
	f, _ := os.OpenFile("1.5w_npm_meta_2021.json", os.O_RDONLY, 0777)
	metas := []SampleStr{}
	decoder := json.NewDecoder(f)
	fmt.Println("test")
	decoder.Decode(&metas)
	fmt.Println(metas)
	f.Close()
}
