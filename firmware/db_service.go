package main

import (
	"database/sql"
	"fmt"
	"time"
)

func insertTable(db sql.DB, firmware Firmware, iphoneName string, tableName string)(err error) {
	fmt.Printf("\"insertTable\": %v\n", "insertTable")
	fmt.Printf("firmware: %v\n", firmware)
	var sql = fmt.Sprintf("INSERT INTO %s(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES (?,?,?,?,?,?,?,?,?,?,?)", 
	tableName,
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
	column_signed)
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	defer stmt.Close()

	var t1 = parseDateStr(firmware.Releasedate)
	fmt.Printf("t1: %v\n", t1)
	var t2 = parseDateStr(firmware.Uploaddate)
	fmt.Printf("t2: %v\n", t2)
	_,er := stmt.Exec(
		firmware.Name,
		firmware.Identifier,
		firmware.Version,
		firmware.Buildid,
		firmware.Sha1Sum,
		firmware.Md5Sum,
		firmware.Filesize,
		firmware.URL,
		t1,
		t2,
		firmware.Signed)
	if er != nil {
		fmt.Printf("er: %v\n", er)
		panic(er)
	}
	fmt.Printf("\"insert success\": %v\n", "insert success")
	return nil
}

func clearFirmware(db sql.DB, productCode string, tableName string) {
	fmt.Printf("\"clearFirmware\": %v\n", "clearFirmware")
	re, err := db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s = "%s"`, tableName, column_identifier, productCode))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	rowsAffeted, er := re.RowsAffected()

	if er != nil {
		fmt.Printf("er: %v\n", er)
	}

	fmt.Printf("rowsAffeted: %v\n", rowsAffeted)
}

///一段时间后重新请求 ipsw 更新
func tickerUpdateFirmware(db sql.DB) {
	// already distinct
	var codes = mapFirmwaresCodes(queryAllFirmware(db))
	for i:=0;i<len(codes);i++ {
		info,er := getIpswInfo(codes[i])
		if er != nil {
			fmt.Printf("er: %v\n", er)
			continue
		}
		time.Sleep(3 * time.Second)
		useInfo := mapToSignedIpsw(info)
		updateFirmware(db, codes[i], useInfo.Name, info.Firmwares, table_firmware)
	}
}


/// 立即更新，接口
func updateFirmware(db sql.DB, productCode string, iphoneName string, firmwares []Firmware, tableName string) {
	fmt.Printf("\"updateFirmware\": %v\n", "updateFirmware")
	clearFirmware(db, productCode, table_firmware)
	for i:= 0; i< len(firmwares); i++ {
		insertTable(db, firmwares[i], iphoneName, tableName)
	}
}

func queryAllFirmware(db sql.DB) ([] Firmware){
	rows,err := db.Query(fmt.Sprintf("SELECT * FROM %s", table_firmware))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return mapRowsToFirmwares(rows)
}

func queryFirmware(db sql.DB, productCode string, tableName string) ([] Firmware){
	var sql = fmt.Sprintf(`SELECT * FROM %s WHERE %s = "%s" AND %s = 1 ORDER BY %s DESC`, tableName, column_identifier, productCode, column_signed, column_releaseDate)
	fmt.Printf("sql: %v\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return mapRowsToFirmwares(rows)
}

/// 获取两个表中的相应 signed firmware 数据, beta 在前面
func generateFirmware(db sql.DB, productCode string) ([] Firmware) {
	betaFirmwares := queryFirmware(db, productCode, table_firmware_beta)
	normalFirmwares := queryFirmware(db, productCode, table_firmware)

	var res [] Firmware
	for i:=0;i<len(betaFirmwares);i++ {
		var t = betaFirmwares[i]
		t.Beta = true
		res = append(res, t)
	}

	for i:=0;i<len(normalFirmwares);i++ {
		var t = normalFirmwares[i]
		t.Beta = false
		res = append(res, t)
	}

	return res
 }

func parseDateStr(dateStr string) time.Time {
	const layout = "2006-01-02T15:04:05Z"
    t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Printf("parse date with err: %v\n", err)
		panic(err)
	}
	return t
}