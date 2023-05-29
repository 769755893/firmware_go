package db

import (
	"database/sql"
	"fmt"
	"main/model"
	"main/util"
	"time"
)

func insertTable(firmware model.Firmware, iphoneName string, tableName string) (err error) {
	fmt.Printf("\"insertTable\": %v\n", "insertTable")
	fmt.Printf("firmware: %v\n", firmware)
	var sqlStr = fmt.Sprintf("INSERT INTO %s(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		tableName,
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
		columnSigned)
	stmt, err := DataBase.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			print(err)
		}
	}(stmt)

	var t1 = parseDateStr(firmware.ReleaseDate)
	fmt.Printf("t1: %v\n", t1)
	var t2 = parseDateStr(firmware.UploadDate)
	fmt.Printf("t2: %v\n", t2)
	_, er := stmt.Exec(
		firmware.Name,
		firmware.Identifier,
		firmware.Version,
		firmware.BuildId,
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

func clearFirmware(productCode string, tableName string) {
	fmt.Printf("\"clearFirmware\": %v\n", "clearFirmware")
	re, err := DataBase.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s = "%s"`, tableName, columnIdentifier, productCode))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	rowsAffected, er := re.RowsAffected()

	if er != nil {
		fmt.Printf("er: %v\n", er)
	}

	fmt.Printf("rowsAffected: %v\n", rowsAffected)
}

// UpdateFirmware / 立即更新，接口
func UpdateFirmware(productCode string, iphoneName string, firmwares []model.Firmware, tableName string) {
	fmt.Printf("\"updateFirmware\": %v\n", "updateFirmware")
	clearFirmware(productCode, TableFirmware)
	for i := 0; i < len(firmwares); i++ {
		err := insertTable(firmwares[i], iphoneName, tableName)
		if err != nil {
			panic(err)
		}
	}
}

func QueryAllFirmware() []model.Firmware {
	rows, err := DataBase.Query(fmt.Sprintf("SELECT * FROM %s", TableFirmware))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return util.MapRowsToFirmwares(rows)
}

func QueryFirmware(productCode string, tableName string) []model.Firmware {
	var sqlStr = fmt.Sprintf(`SELECT * FROM %s WHERE %s = "%s" AND %s = 1 ORDER BY %s DESC`, tableName, columnIdentifier, productCode, columnSigned, columnReleaseDate)
	fmt.Printf("sqlStr: %v\n", sqlStr)
	rows, err := DataBase.Query(sqlStr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return util.MapRowsToFirmwares(rows)
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
