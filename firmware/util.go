package main

import (
	"database/sql"
	"fmt"
)


func mapToSignedFirmware(ipsw *IPSWFirm) (firmware [] Firmware) {
	fimwares := ipsw.Firmwares
	var ret = make([] Firmware, 1)
	for i := 0; i < len(fimwares); i++ {
		if (fimwares[i].Signed) {
			ret = append(ret, fimwares[i])
		}
	}
	return ret;
}

func mapRowsToFirmwares(rows *sql.Rows) ([] Firmware) {
	var firmwares [] Firmware
	for rows.Next() {
		var firmware Firmware
		err := rows.Scan(&firmware.Name, &firmware.Identifier, &firmware.Buildid, &firmware.Sha1Sum, &firmware.Md5Sum, &firmware.Filesize, &firmware.URL, &firmware.Releasedate, &firmware.Uploaddate, &firmware.Signed)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		firmwares = append(firmwares, firmware)
	}
	return firmwares
}

func mapFirmwaresCodes(firmwares [] Firmware) ([] string) {
	var ret [] string
	for i:=0;i<len(firmwares);i++ {
		ret = append(ret, firmwares[i].Identifier)
	}
	return uniqueMap(ret)
}

func uniqueMap(code [] string)([] string) {
	result := make([]string, 0, len(code))
    set := make(map[string]struct{})
    for _, n := range code {
        if _, ok := set[n]; !ok {
            set[n] = struct{}{}
            result = append(result, n)
        }
    }
    return result
}