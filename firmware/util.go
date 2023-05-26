package main

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode"
)

func mapToSignedIpsw(ipsw *IPSWFirm) *IPSWFirm {
	var firmwares [] Firmware
	firmwares = mapToSignedFirmware(ipsw)
	return &IPSWFirm{
		Name: ipsw.Name,
		Identifier: ipsw.Identifier,
		Firmwares: firmwares,
		Boards: ipsw.Boards,
		Boardconfig: ipsw.Boardconfig,
		Platform: ipsw.Platform,
		Cpid: ipsw.Cpid,
		Bdid: ipsw.Bdid,
	}
}

func mapToSignedFirmware(ipsw *IPSWFirm)[] Firmware {
	fimwares := ipsw.Firmwares
	var name = ipsw.Name
	var ret [] Firmware
	for i := 0; i < len(fimwares); i++ {
		if (fimwares[i].Signed) {
			var t = fimwares[i]
			t.Name = name
			ret = append(ret, t)
		}
	}
	return ret;
}

func mapRowsToFirmwares(rows *sql.Rows) ([] Firmware) {
	fmt.Printf("\"start map Rows to firmwares\": %v\n", "start map Rows to firmwares")
	var firmwares [] Firmware
	for rows.Next() {
		var firmware Firmware
		var id int
		err := rows.Scan(
			&id,
			&firmware.Name,
			&firmware.Identifier,
			&firmware.Version,
			&firmware.Buildid,
			&firmware.Sha1Sum,
			&firmware.Md5Sum,
			&firmware.Filesize,
			&firmware.URL,
			&firmware.Releasedate,
			&firmware.Uploaddate,
			&firmware.Signed)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		firmwares = append(firmwares, firmware)
	}
	fmt.Printf("\"map the end\": %v\n", "map the end")
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

func replaceLetter(str string, n int) string {
	var done = false
	return strings.Map(func (r rune) rune {
		if unicode.IsLetter(r) && !done {
			if (n == 1) {
				done = true
			}
			return ' '
		}
		return r
	}, str)
}