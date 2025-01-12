package main

import (
	"fmt"
	"rest-api-go/internal/config"
	"rest-api-go/internal/repository"
	"rest-api-go/pkg"
	"rest-api-go/pkg/postgres"
	"rest-api-go/pkg/redis"
	"rest-api-go/web"
)

func main() {
	connectionStringPostgres := fmt.Sprintf(`host=%s
	dbname=%s
	user=%s
	password=%s
	port=%d
	sslmode=disable`,
		config.Cfg.Database.Host,
		config.Cfg.Database.DBName,
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.Port,
	)

	connectionStringRedis := fmt.Sprintf("redis://%s@%s:%s/%d",
		config.Cfg.RedisDatabase.Password,
		config.Cfg.RedisDatabase.Host,
		config.Cfg.RedisDatabase.Port,
		config.Cfg.RedisDatabase.RedisDB,
	)

	postgresConnection := postgres.ConnectionPostgres{connectionStringPostgres}
	con := pkg.DBConnection{Db: postgresConnection}
	DBPostgres := con.DBConnect()

	redisConnection := redis.ConnectionRedis{connectionStringRedis}
	con1 := pkg.DBConnection{Db: redisConnection}
	con1.DBConnect()

	userRepo := repository.NewUserRepository(DBPostgres)

	//Initialize setup for Org1
	cryptoPath := "../../test-network/organizations/peerOrganizations/org1.example.com"
	orgConfig := web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
		KeyPath:      cryptoPath + "/users/User1@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
		UserRepo:     *userRepo,
	}

	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		fmt.Println("Error initializing setup for Org1: ", err)
	}
	web.Serve(web.OrgSetup(*orgSetup))
}
