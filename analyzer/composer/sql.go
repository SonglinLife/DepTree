package composer

import (
	"context"
	"fmt"
	"opendep/mongo"
)

/**
 * @description: 根据组件名称和version查询对应的 composer.json 信息
 * @param {*} name
 * @param {string} version
 * @return {*}
 */
func ComposerRecords(name, version string) (*Records, bool) {
	record := NewRecords()
	sql := `
	[
		{"$match":{"name":"%s", "version":"%s"}},
		{"$project":{"name":1, "version":1, "require":"$detail.require", "license":"$detail.license"}},
		{"$limit":1}
	]
	`
	sql = fmt.Sprintf(sql, name, version)
	cursor, err := mongo.Query(sql, "composer_records")
	fl := false

	if err != nil {
		return record, fl
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		cursor.Decode(record)
		fl = true
		break
	}
	return record, fl
}

/**
 * @description: 根据组件名称得到所有的Meta数据
 * @param {string} name
 * @return {*}
 */
func ComposerMetas(name string) []Metas {
	meta := []Metas{}
	sql := `
		[
			{"$match":{"name":"%s"}},
			{"$unwind":"$releases"},
			{"$project":{"_id":0, "time":"$releases.time", "name":1, "version":"$releases.version"}}
		]
	`
	sql = fmt.Sprintf(sql, name)

	cursor, err := mongo.Query(sql, "composer_metas")
	if err != nil {
		return meta
	}
	defer cursor.Close(context.Background())

	cursor.All(context.Background(), &meta)

	return meta
}

func VersionTime(name, version string) {

}
