package sample

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"opendep/analyzer/npm"
	"opendep/depstr"
	"opendep/mongo"
	"opendep/mttu"
	"os"
	"time"

	"github.com/Jeffail/tunny"
)

type SampleStr struct {
	Name     string    `bson:"name" json:"name"`
	Releases []Release `bson:"releases" json:"releases"`
}
type Release struct {
	Version string    `bson:"version" json:"version"`
	Time    time.Time `bson:"time" json:"time"`
}

/**
 * @description: 该函数仅做示例
 * @param {*}
 * @return {*}
 */
func SampleMetas(collection string) {

	sql := `
		[
			{"$project":{
				"name":1,
				"_id":0,
				"releases":{
				   "$filter":{
					"input":"$releases",
					"as":"release",
					"cond":{
						"$and":[
							{"$gte": ["$$release.time",ISODate("2021-01-01T00:00:00.000Z")]},
							{"$lt":  ["$$release.time",ISODate("2022-01-01T00:00:00.000Z")]}
						]
					}
				}
				}
			}},
			{"$match":{"releases": {"$ne": []}}},
			{"$sample":{"size":15000}}
		]	
	`
	metas := []SampleStr{}
	fName := fmt.Sprintf("1.5w_%v_2021.json", collection)
	f, _ := os.OpenFile(fName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	encoder := json.NewEncoder(f)
	cursor, err := mongo.Query(sql, collection)
	if err != nil {
		log.Println(err)
	}
	cursor.All(context.Background(), &metas)
	encoder.Encode(metas)
	f.Close()
}

type External struct {
	Name        string    `json:"name"`
	Version     string    `json:"verion"`
	Time        time.Time `json:"time"`
	Externaldep float32   `json:"externaldep"`
}

func SampleExternal() {

	f, _ := os.OpenFile("1.5w_npm_meta_2021.json", os.O_RDONLY, 0777)

	metas := []SampleStr{}

	decoder := json.NewDecoder(f)
	decoder.Decode(&metas)

	ef, _ := os.OpenFile("external.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)

	encoder := json.NewEncoder(ef)

	// e := []External{}
	pool := tunny.NewFunc(3, func(i interface{}) interface{} {
		meta := i.(SampleStr)
		n := rand.Intn(len(meta.Releases))
		root := depstr.NewDepTree(nil, meta.Name, meta.Releases[n].Version)
		npm.BuildDepTree(root)
		external := root.External()
		ex := External{
			Name:        meta.Name,
			Version:     meta.Releases[n].Version,
			Time:        meta.Releases[n].Time,
			Externaldep: external,
		}

		encoder.Encode(ex)

		return nil
	})
	for _, meta := range metas {
		pool.Process(meta)
	}
	// encoder.Encode(e)
}

type MTTU struct {
	Name    string    `json:"name"`
	Version string    `json:"verion"`
	Time    time.Time `json:"time"`
	Mttu    float32   `json:"mttu"`
}

func SampleMTTU() {
	rand.Seed(time.Now().UnixMilli())
	f, _ := os.OpenFile("1.5w_npm_meta_2021.json", os.O_RDONLY, 0777)

	metas := []SampleStr{}

	decoder := json.NewDecoder(f)
	decoder.Decode(&metas)

	ef, _ := os.OpenFile("mttu.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)

	encoder := json.NewEncoder(ef)

	// e := []External{}
	pool := tunny.NewFunc(3, func(i interface{}) interface{} {
		meta := i.(SampleStr)
		n := rand.Intn(len(meta.Releases))
		root := depstr.NewDepTree(nil, meta.Name, meta.Releases[n].Version)
		root.Time = meta.Releases[n].Time
		mttu.BuildMttuTree(root)
		mttu := root.MTTU()
		ex := MTTU{
			Name:    meta.Name,
			Version: meta.Releases[n].Version,
			Time:    meta.Releases[n].Time,
			Mttu:    mttu,
		}

		encoder.Encode(ex)
		return nil
	})
	for _, meta := range metas {
		pool.Process(meta)
	}
}
