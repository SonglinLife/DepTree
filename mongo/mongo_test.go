
package mongo

import (
	"context"
	"log"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

/**
 * @description: 测试mongoDB连接时候正常
 * @param {*testing.T} t
 * @return {*}
 */
func TestConnect(t *testing.T) {

	client, err := InitDb()

	if err != nil {
		log.Fatal(err)
	}
	//检查是否ping通
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// 列出database名称
	result, err := client.ListDatabaseNames(
		context.TODO(),
		bson.M{})
	log := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	if err != nil {
		log.Fatal(err)
	}
	for _, db := range result {
		log.Println("database: " + db)
	}
	// 列出database admin中的所有collection
	db := client.Database("admin")
	cols, _ := db.ListCollectionNames(context.TODO(), bson.M{})

	for _, col := range cols {
		log.Println("col: " + col)
	}

}
