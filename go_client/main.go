package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

var (
	dbUser = "alice"
	dbName = "testdb"
	cacert = flag.String("cacert", "../certs/root-ca.crt", "RootCA")
)

func main() {

	flag.Parse()

	caCert, err := os.ReadFile(*cacert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	r, err := tls.LoadX509KeyPair("../certs/pg-client.crt", "../certs/pg-client.key")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("user=%s database=%s sslmode=verify-ca", dbUser, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Printf("cloudsqlconn.NewDialer: %v", err)
		os.Exit(1)
	}

	config.Password = "somepassword"
	config.Host = "127.0.0.1"
	//config.Port = 5432 // direct
	config.Port = 15432 // envoy
	config.TLSConfig = &tls.Config{
		//serverName: "postgres.domain.com", // direct
		ServerName:   "envoy.domain.com", // envoy
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{r},
	}

	sslKeyLogfile := os.Getenv("SSLKEYLOGFILE")
	if sslKeyLogfile != "" {
		var w *os.File
		w, err := os.OpenFile(sslKeyLogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatalf("Could not create keylogger: ", err)
		}
		config.TLSConfig.KeyLogWriter = w
	}

	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Printf("cloudsqlconn.NewDialer: %v", err)
		os.Exit(1)
	}

	err = dbPool.Ping()
	if err != nil {
		log.Printf("cloudsqlconn.ping: %v", err)
		os.Exit(1)
	}
	log.Println("Done")

}
