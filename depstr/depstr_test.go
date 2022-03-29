package depstr

import (
	"context"
	"fmt"
	"log"
	"opendep/mongo"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

/**
 * @description: 用于测试是否可以正常编码
 * @param {*testing.T} t
 * @return {*}
 */
func TestDepstruct(t *testing.T) {
	client, err := mongo.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database("admin").Collection("npm_records")

	filter := bson.M{}

	var result Records

	coll.FindOne(context.TODO(), filter).Decode(&result)

	fmt.Println(result)

}

/**
 * @description: 测试管道函数
 * @param {*testing.T} t
 * @return {*}
 */
func TestAggregate(t *testing.T) {

	sqlstr := `[{"$match":{"name":"response-json-formatter"}},
	{"$project":{"name":1, "version":1}},
	{"$limit":100}
	]`

	client, _ := mongo.InitDB()
	npm := client.Database("admin").Collection("npm_records")
	cursor, _ := npm.Aggregate(context.TODO(), mongo.Sqlstr2Bson(sqlstr))
	var result []Records
	cursor.All(context.TODO(), &result)

	fmt.Println(result)

}

/**
 * @description: 测试能否获得meta数据
 * @param {*testing.T} t
 * @return {*}
 */
func TestMetas(t *testing.T) {
	sql := `
		[
			{"$match":{"name":"response-json-formatter"}},
			{"$unwind":"$releases"},
			{"$project":{"_id":0, "time":"$releases.time", "name":1, "version":"$releases.version"}},
		]
	`

	client, _ := mongo.InitDB()
	npm_metas := client.Database("admin").Collection("npm_metas")
	var result []Records

	cursor, _ := npm_metas.Aggregate(context.Background(), mongo.Sqlstr2Bson(sql))
	cursor.All(context.TODO(), &result)
	fmt.Println(result)
}

/**
 * @description: 测试能否得到dependencies数据
 * @param {*testing.T} t
 * @return {*}
 */
func TestGetdep(t *testing.T) {
	sql := `
	[
		{"$match":{"name":"field-descriptions", "version":"1.0.7"}},
		{"$project":{"name":1, "version":1, "devDependencies":"$detail.devDependencies", "Dependencies":"$detail.dependencies"}},
		{"$limit":1}
	]
	`
	record := NewRecords()


	cursor, _ := mongo.Query(sql,"npm_records")
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		cursor.Decode(record)
		break
	}
	fmt.Println(record)
	
}

/**
 * @description: 测试连表查询
 * @param {*testing.T} t
 * @return {*}
 */
func TestUnion(t *testing.T) {
	sql := `
	[
		{$match:{"name":"field-descriptions", version:"1.0.8"}},
		{$project:{"_id":0, "devDependencies":"$detail.devDependencies", "name":1, "version":1}},
		{$unionWith: {
			"coll": "npm_metas",
			"pipeline":[
				{"$match":{"name":"field-descriptions"}},
				{"$unwind":"$releases"},
				{"$project":{"_id":0, "time":"$releases.time","tversion":"$releases.version"}}
			]
		}}
		
	]
	`
	mongo.Sqlstr2Bson(sql)
	// client, _ := mongo.InitDB()
	// npm_metas := client.Database("admin").Collection("npm_records")
	// var result []Records

	// cursor, _ := npm_metas.Aggregate(context.Background(), mongo.Sqlstr2Bson(sql))
	// cursor.All(context.TODO(), &result)
	// fmt.Println(result)

}
