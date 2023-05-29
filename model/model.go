package model

import (
	JsonIter "github.com/json-iterator/go"
)

type IPSWFirm struct {
	Name        string     `json:"name"`
	Identifier  string     `json:"identifier"`
	Firmwares   []Firmware `json:"firmwares"`
	Boards      []Board    `json:"boards"`
	BoardConfig string     `json:"boardconfig"`
	Platform    string     `json:"platform"`
	CpId        int64      `json:"cpid"`
	BdId        int64      `json:"bdid"`
}

type Board struct {
	BoardConfig string `json:"boardconfig"`
	Platform    string `json:"platform"`
	CpId        int64  `json:"cpid"`
	BdId        int64  `json:"bdid"`
}

type Firmware struct {
	Name        string `json:name`
	Identifier  string `json:"identifier"`
	Version     string `json:"version"`
	BuildId     string `json:"buildid"`
	Sha1Sum     string `json:"sha1sum"`
	Md5Sum      string `json:"md5sum"`
	Sha256Sum   string `json:"sha256sum"`
	Filesize    int64  `json:"filesize"`
	URL         string `json:"url"`
	ReleaseDate string `json:"releasedate"`
	UploadDate  string `json:"uploaddate"`
	Signed      bool   `json:"signed"`
	Beta        bool   `json:"beta"`
}

type Identifier string

var json = JsonIter.ConfigCompatibleWithStandardLibrary

const (
	IPhone81 Identifier = "iPhone8,1"
)

func DecodeIpsw(data []byte) (*IPSWFirm, error) {
	var r IPSWFirm
	err := json.Unmarshal(data, &r)
	return &r, err
}

func EncodeIpsw(r *IPSWFirm) ([]byte, error) {
	return json.Marshal(r)
}
