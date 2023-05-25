package server

import (
	"database/sql"
	"fmt"
)

const ipswFirmwareTable = "IpswFirmware"
const ipswFirmwareBetaTable = "IpswFirmwareBeta"
const column_version = "Version"
const column_BuildId = "Buildid"
const column_shaSum = "Sha1Sum"
const column_md5Sum = "Md5Sum"
const column_fileSize = "Filesize"
const column_url = "URL"
const column_releaseDate = "Releasedate"
const column_uploaddate = "Uploaddate"
const column_signed = "Signed"


func createIpswTable(db sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		 id integer PRIMARY KEY,
		 %s varchar(255),
		%s varchar(255),
		%s varchar(255),
		%s varchar(255),
		%s varchar(255),
		%s integer,
		%s string,
		%s timestamp,
		%s timestamp,
		%s bool
	  );
		`, ipswFirmwareTable,
		column_version,
		column_BuildId,
		column_shaSum,
		column_md5Sum,
		column_shaSum,
		column_fileSize,
		column_fileSize,
		column_url,
		column_releaseDate,
		column_uploaddate,
		column_signed,
		))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\"create firmware table success\": %v\n", "create firmware table success")
}

func createBetaTable(db sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		 id integer PRIMARY KEY,
		 %s varchar(255),
		%s varchar(255),
		%s varchar(255),
		%s varchar(255),
		%s varchar(255),
		%s integer,
		%s string,
		%s timestamp,
		%s timestamp,
		%s bool
	  );
		`, ipswFirmwareBetaTable,
		column_version,
		column_BuildId,
		column_shaSum,
		column_md5Sum,
		column_shaSum,
		column_fileSize,
		column_fileSize,
		column_url,
		column_releaseDate,
		column_uploaddate,
		column_signed,
		))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\"create firmware table success\": %v\n", "create firmware table success")
}