package mongo

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/elliotchance/sshtunnel"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh"
)

func initTunnel() {
	tunnel := sshtunnel.NewSSHTunnel(
		// ip地址
		"21xlabStudent@139.196.6.95:22",
		// 密码
		ssh.Password("xlabUser666/"),
		// 数据库地址
		"dds-uf61fd4fb14bd3541744-pub.mongodb.rds.aliyuncs.com:3717",
		// 本地绑定端口
		"3733",
	)
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	tunnel.Log = logger
	go tunnel.Start()
}

func InitDb() (*mongo.Client, error) {

	initTunnel()
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      "xlab",
		Password:      "Xlab2021!",
		AuthSource:    "admin",
		PasswordSet:   true,
	}

	clientOpts := options.Client().ApplyURI("mongodb://localhost:3733/?connect=direct").
		SetAuth(credential)
	
	//连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	return client, err
}
