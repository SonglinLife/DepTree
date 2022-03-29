package main

import (
	"encoding/json"
	"fmt"
	"log"
	"opendep/analyzer/npm"
	"opendep/depstr"
	"opendep/mttu"
	"os"
)

/**
 * @description: 用于计算Npm的MTTU
 * @param {*} name
 * @param {string} version
 * @return {*}
 */
func NpmMttu(name, version string) {
	ti := npm.VersionTime(name, version)
	if ti == nil {
		log.Fatalln("no such package!")
	}
	root := depstr.NewDepTree(nil, name, version)
	root.Time = *ti
	mttu.BuildMttuTree(root)
	mttu := root.MTTU()
	fmt.Printf("%s@%s mttu:%f\n", name, version, mttu)
}

/**
 * @description: 用于计算Npm的外部依赖度
 * @param {*} name
 * @param {string} version
 * @return {*}
 */
func NpmExternal(name, version string) {
	ti := npm.VersionTime(name, version)
	if ti == nil {
		log.Fatalln("no such package!")
	}
	root := depstr.NewDepTree(nil, name, version)
	root.Time = *ti
	npm.BuildDepTree(root)
	external := root.External()
	fmt.Printf("%s@%s external:%f\n", name, version, external)

}

/**
 * @description: 拉取该组件的所有历史版本号和时间
 * @param {string} name
 * @return {*}
 */
func NpmAllVersion(name string) {
	meta := npm.NpmMetas(name)
	s := fmt.Sprintf("%v_all_version.json", name)
	f, err := os.OpenFile(s, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err!=nil{
		log.Fatal(err)
	}
	encoder := json.NewEncoder(f)
	encoder.Encode(meta)
	f.Close()
}


/**
 * @description: 保存该依赖树
 * @param {*} name
 * @param {string} version
 * @return {*}
 */
func NpmSaveTree(name, version string){
	ti := npm.VersionTime(name, version)
	if ti == nil {
		log.Fatalln("no such package!")
	}
	root := depstr.NewDepTree(nil, name, version)
	root.Time = *ti
	npm.BuildDepTree(root)
	root.LevelOrder()
	root.PrintTree()
}

