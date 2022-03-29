package mongo

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/elliotchance/sshtunnel"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh"
	"gopkg.in/mgo.v2/bson"
)

type cfg struct {
	TunnelIp       string `json:"tunnelIp"`
	TunnelPassword string `json:"tunnelPasswd"`
	DBAddr         string `json:"dbAddr"`
	DBUser         string `json:"dbUser"`
	DBPassword     string `json:"dbPassword"`
}

func loadConfig() *cfg {
	c := cfg{}
	var f, err = os.Open("config.json")
	if err != nil {
		log.Fatal("No Config Json!")
	}
	decoder := json.NewDecoder(f)
	decoder.Decode(&c)
	return &c
}

var Client *mongo.Client

func initTunnel(ip, passwd, dbAddr string) {
	tunnel := sshtunnel.NewSSHTunnel(
		// ip地址
		ip,
		// 密码
		ssh.Password(passwd),
		// 数据库地址
		dbAddr,
		// 本地绑定端口
		"33777",
	)
	go tunnel.Start()
}

func InitDB() (*mongo.Client, error) {
	if Client != nil {
		return Client, nil
	}
	c := loadConfig()
	initTunnel(c.TunnelIp, c.TunnelPassword, c.DBAddr)
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      c.DBUser,
		Password:      c.DBPassword,
		AuthSource:    "admin",
		PasswordSet:   true,
	}

	// 限制连接池最大为200
	clientOpts := options.Client().ApplyURI("mongodb://localhost:33777/?connect=direct").
		SetAuth(credential).SetMaxPoolSize(200)

	//连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err == nil {
		Client = client
	}
	return client, err
}

/**
 * @description: 将查询mongodb的语句编码为bson
 * @param {string} sql
 * @return {*}
 */
func Sqlstr2Bson(sql string) interface{} {
	var bdoc interface{}
	bson.UnmarshalJSON([]byte(sql), &bdoc)
	return bdoc
}


/**
 * @description: 将查询包装为一个函数 
 * @param {string} sqlstr
 * @param {string} coll
 * @return {*}
 */
func Query(sqlstr string, coll string) (*mongo.Cursor, error) {
	sql := Sqlstr2Bson(sqlstr)
	client, _ := InitDB()
	col := client.Database("admin").Collection(coll)
	return col.Aggregate(context.Background(), sql)
}
