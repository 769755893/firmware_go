package main

import (
	"database/sql"
	"fmt"
)

const table_firmware = "IpswFirmware"
const table_firmware_beta = "IpswFirmwareBeta"

const column_iphoneName = "name"//1
const column_identifier = "identifier"//2
const column_version = "version" //3
const column_BuildId = "buildId" //4
const column_shaSum = "sha1Sum" //5
const column_md5Sum = "md5Sum" //6
const column_fileSize = "fileSize" //7
const column_url = "url" //8
const column_releaseDate = "releaseDate" //9
const column_uploaddate = "uploaddate" //10
const column_signed = "available" //11

func createIpswTable(db sql.DB, tableName string) (er error){
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
		column_iphoneName,
		column_identifier,
		column_version,
		column_BuildId,
		column_shaSum,
		column_md5Sum,
		column_fileSize,
		column_url,
		column_releaseDate,
		column_uploaddate,
		column_signed,
		))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return er
	}

	fmt.Printf("\"create firmware table success\": %v\n", "create firmware table success")
	return nil
}

func initTable(db sql.DB)(er error) {
	er0 := createIpswTable(db, table_firmware)
	if er0 != nil {
		fmt.Printf("failed to create ipsw table: %v\n", er0)
	}

	er1 := createIpswTable(db, table_firmware_beta)
	if er1 != nil {
		fmt.Printf("failed to create ipsw beta table: %v\n", er1)
	}

	fmt.Printf("\"success to create all table\": %v\n", "success to create all table")

	fmt.Printf("\"now the table status\": %v\n", "now the table status")

	return nil
}

func deleteAllTable(db sql.DB) {
	_, e0 := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table_firmware))
	_, e1 := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table_firmware_beta))
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