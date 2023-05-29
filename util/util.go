package util

import (
	"database/sql"
	"fmt"
	"main/model"
	"net/http"
	"strings"
	"unicode"
)

func ByteString(str string) []byte {
	return []byte(str)
}

func HttpWriter(w http.ResponseWriter, code int, bytes []byte) {
	w.WriteHeader(code)
	_, err := w.Write(bytes)
	if err != nil {
		fmt.Printf("server err: %s", err)
		return
	}
}

func MapToSignedIpsw(ipsw *model.IPSWFirm) *model.IPSWFirm {
	var firmwares []model.Firmware
	firmwares = MapToSignedFirmware(ipsw)
	return &model.IPSWFirm{
		Name:        ipsw.Name,
		Identifier:  ipsw.Identifier,
		Firmwares:   firmwares,
		Boards:      ipsw.Boards,
		BoardConfig: ipsw.BoardConfig,
		Platform:    ipsw.Platform,
		CpId:        ipsw.CpId,
		BdId:        ipsw.BdId,
	}
}

func MapToSignedFirmware(ipsw *model.IPSWFirm) []model.Firmware {
	firmwares := ipsw.Firmwares
	var name = ipsw.Name
	var ret []model.Firmware
	for i := 0; i < len(firmwares); i++ {
		if firmwares[i].Signed {
			var t = firmwares[i]
			t.Name = name
			ret = append(ret, t)
		}
	}
	return ret
}

func MapRowsToFirmwares(rows *sql.Rows) []model.Firmware {
	fmt.Printf("\"start map Rows to firmwares\": %v\n", "start map Rows to firmwares")
	var firmwares []model.Firmware
	for rows.Next() {
		var firmware model.Firmware
		var id int
		err := rows.Scan(
			&id,
			&firmware.Name,
			&firmware.Identifier,
			&firmware.Version,
			&firmware.BuildId,
			&firmware.Sha1Sum,
			&firmware.Md5Sum,
			&firmware.Filesize,
			&firmware.URL,
			&firmware.ReleaseDate,
			&firmware.UploadDate,
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

func MapFirmwaresCodes(firmwares []model.Firmware) []string {
	var ret []string
	for i := 0; i < len(firmwares); i++ {
		ret = append(ret, firmwares[i].Identifier)
	}
	return uniqueMap(ret)
}

func uniqueMap(code []string) []string {
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
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) && !done {
			if n == 1 {
				done = true
			}
			return ' '
		}
		return r
	}, str)
}
