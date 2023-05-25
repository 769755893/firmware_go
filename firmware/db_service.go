package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

func insertTable(db sql.DB, firmware Firmware, iphoneName string, tableName string)(err error) {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES (?,?,?,?,?,?,?,?,?)", tableName, column_iphoneName ,column_identifier,column_version, column_BuildId, column_shaSum, column_md5Sum, column_fileSize, column_url, column_releaseDate, column_uploaddate, column_signed))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	defer stmt.Close()

	var t1 = parseDateStr(firmware.Releasedate)
	var t2 = parseDateStr(firmware.Uploaddate)
	_,er := stmt.Exec(firmware.Version, firmware.Buildid, firmware.Sha1Sum, firmware.Md5Sum, firmware.Filesize, firmware.URL, t1, t2, firmware.Signed)
	if er != nil {
		fmt.Printf("insert failed with er: %v\n", er)
		return er
	}
	fmt.Printf("\"insert success\": %v\n", "insert success")
	return nil
}

func clearFirmware(db sql.DB, productCode string, tableName string) {
	re, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = %s", tableName, column_identifier, productCode))
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
		info,er :=getIpswInfo(codes[i])
		if er != nil {
			fmt.Printf("er: %v\n", er)
			continue
		}
		time.Sleep(3 * time.Second)
		updateFirmware(db, codes[i], info.Name, info.Firmwares, table_firmware)
	}
}


/// 立即更新，接口
func updateFirmware(db sql.DB, productCode string, iphoneName string, firmwares []Firmware, tableName string) {
	clearFirmware(db, productCode, table_firmware)
	for i:= 0; i< len(firmwares); i++ {
		insertTable(db, firmwares[i], iphoneName, tableName)
	}
}

func queryAllFirmware(db sql.DB) ([] Firmware){
	rows,err := db.Query(fmt.Sprintf("SELECT * FROM %s", table_firmware))
	if err != nil {
		panic(err)
	}
	return mapRowsToFirmwares(rows)
}

func queryFirmware(db sql.DB, productCode string, tableName string) ([] Firmware){
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s = %s AND %s = 1 ORDER BY %s DESC", tableName, column_identifier, productCode, column_signed, column_releaseDate))
	if err != nil {
		panic(err)
	}
	return mapRowsToFirmwares(rows)
}

/// 获取两个表中的相应 signed firmware 数据, beta 在前面
func generateFirmware(db sql.DB, productCode string) ([] Firmware) {
	betaFirmwares := queryFirmware(db, productCode, table_firmware_beta)
	normalFirmwares := queryFirmware(db, productCode, table_firmware)

	var res [] Firmware
	for i:=0;i<len(betaFirmwares);i++ {
		t := reflect.ValueOf(betaFirmwares[i])
		t.FieldByName("Beta").SetBool(true)

		var firmware = t.Interface().(Firmware)
		res = append(res, firmware)
	}

	for i:=0;i<len(normalFirmwares);i++ {
		t := reflect.ValueOf(normalFirmwares[i])
		t.FieldByName("Beta").SetBool(false)

		var firmware = t.Interface().(Firmware)
		res = append(res, firmware)
	}

	return res
 }


/// TODO 设置超时时间自动更新 表中的 firmware
func updateOutDateFirmware() {

}

func parseDateStr(dateStr string)(date *time.Time) {
	d, e := time.Parse("2023-05-23 10:01", dateStr)
	if e != nil {
		fmt.Printf("e: %v\n", e)
		return nil
	}

	return &d
}