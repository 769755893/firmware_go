package db

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

const TableFirmware = "IpswFirmware"
const TableFirmwareBeta = "IpswFirmwareBeta"

const columnIphoneName = "name"         //1
const columnIdentifier = "identifier"   //2
const columnVersion = "version"         //3
const columnBuildId = "buildId"         //4
const columnShaSum = "sha1Sum"          //5
const columnMd5sum = "md5Sum"           //6
const columnFilesize = "fileSize"       //7
const columnUrl = "url"                 //8
const columnReleaseDate = "releaseDate" //9
const columnUploadDate = "uploaddate"   //10
const columnSigned = "available"        //11

func createIpswTable(db *sql.DB, tableName string) (er error) {
	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		 id INT AUTO_INCREMENT PRIMARY KEY,
		 %s VARCHAR(50),
		 %s VARCHAR(100),
		 %s VARCHAR(100),
		%s VARCHAR(100),
		%s VARCHAR(200),
		%s VARCHAR(200),
		%s BIGINT,
		%s VARCHAR(500),
		%s DATETIME,
		%s DATETIME,
		%s INT(1)
	  );
		`, tableName,
		columnIphoneName,
		columnIdentifier,
		columnVersion,
		columnBuildId,
		columnShaSum,
		columnMd5sum,
		columnFilesize,
		columnUrl,
		columnReleaseDate,
		columnUploadDate,
		columnSigned,
	))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return er
	}

	fmt.Printf("\"create firmware table success\": %v\n", "create firmware table success")
	return nil
}

func InitTable() (er error) {
	er0 := createIpswTable(DataBase, TableFirmware)
	if er0 != nil {
		fmt.Printf("failed to create ipsw table: %v\n", er0)
	}

	er1 := createIpswTable(DataBase, TableFirmwareBeta)
	if er1 != nil {
		fmt.Printf("failed to create ipsw beta table: %v\n", er1)
	}

	fmt.Printf("\"success to create all table\": %v\n", "success to create all table")

	fmt.Printf("\"now the table status\": %v\n", "now the table status")

	return nil
}

func DeleteAllTable() {
	_, e0 := DataBase.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", TableFirmware))
	_, e1 := DataBase.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", TableFirmwareBeta))
	if e0 != nil {
		fmt.Printf("e0: %v\n", e0)
		return
	}
	if e1 != nil {
		fmt.Printf("e1: %v\n", e1)
		return
	}

	fmt.Printf("\"Drop table Success\": %v\n", "Drop table Success")
}

func OpenDB(driverName string, host string, user string, password string, dataBase string) *sql.DB {
	dataSourceName := user + ":" + password + "@" + "tcp" + "(" + host + ")" + "/" + dataBase + "?charset=utf8"

	fmt.Printf("dataSource %s\n", dataSourceName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetConnMaxIdleTime(time.Minute * 10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		fmt.Println("mysql connect failed")
		panic(err)
	}

	fmt.Println("mysql connect complete")
	return db
}

func LoadDbConfig() {
	exePath, er := os.Getwd()
	if er != nil {
		panic(er)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(exePath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	dataBase := viper.GetString("db.database")
	driver := viper.GetString("db.databaseType")
	DataBase = OpenDB(driver, fmt.Sprintf("%s:%s", host, port), user, password, dataBase)
}

var DataBase *sql.DB
