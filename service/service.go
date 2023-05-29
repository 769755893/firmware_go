package service

import (
	"fmt"
	"main/api"
	"main/db"
	"main/model"
	"main/util"
	"time"
)

// GenerateFirmware / 获取两个表中的相应 signed firmware 数据, beta 在前面
func GenerateFirmware(productCode string) []model.Firmware {
	betaFirmwares := db.QueryFirmware(productCode, db.TableFirmwareBeta)
	normalFirmwares := db.QueryFirmware(productCode, db.TableFirmware)

	var res []model.Firmware
	for i := 0; i < len(betaFirmwares); i++ {
		var t = betaFirmwares[i]
		t.Beta = true
		res = append(res, t)
	}

	for i := 0; i < len(normalFirmwares); i++ {
		var t = normalFirmwares[i]
		t.Beta = false
		res = append(res, t)
	}

	return res
}

func SaveIpswToDB(firmwareCode string, iphoneName string, firmwares []model.Firmware) {
	fmt.Printf("\"saveIpswToDB\": %v\n", "saveIpswToDB")
	fmt.Printf("firmwares: %v\n", firmwares)
	db.UpdateFirmware(firmwareCode, iphoneName, firmwares, db.TableFirmware)
}

// TickerUpdateFirmware /一段时间后重新请求 ipsw 更新
func TickerUpdateFirmware() {
	// already distinct
	var codes = util.MapFirmwaresCodes(db.QueryAllFirmware())
	for i := 0; i < len(codes); i++ {
		info, er := api.GetIpswInfo(codes[i])
		if er != nil {
			fmt.Printf("er: %v\n", er)
			continue
		}
		time.Sleep(3 * time.Second)
		useInfo := util.MapToSignedIpsw(info)
		db.UpdateFirmware(codes[i], useInfo.Name, info.Firmwares, db.TableFirmware)
	}
}

func StartMonitor() {
	ticker := time.Tick(3600 * time.Second)
	fmt.Printf("\"monitor started\": %v\n", "monitor started")

	for range ticker {
		fmt.Printf("\"monitor running update the firmware\": %v\n", "monitor running update the firmware")
		TickerUpdateFirmware()
	}
}

// GetNewIpswFirmware / already sync db data.
func GetNewIpswFirmware(firmwareCode string) (*model.IPSWFirm, error) {
	fmt.Printf("\"getNewIpswFirmware\": %v\n", "getNewIpswFirmware")
	ipsw, err := api.GetIpswInfo(firmwareCode)
	fmt.Printf("get the new ipsw: %v\n", ipsw)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	SaveIpswToDB(firmwareCode, ipsw.Name, util.MapToSignedFirmware(ipsw))
	return util.MapToSignedIpsw(ipsw), nil
}
