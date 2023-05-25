package main

type IPSWFirm struct {
	Name        string     `json:"name"`
	Identifier  string `json:"identifier"`
	Firmwares   []Firmware `json:"firmwares"`
	Boards      []Board    `json:"boards"`
	Boardconfig string     `json:"boardconfig"`
	Platform    string     `json:"platform"`
	Cpid        int64      `json:"cpid"`
	Bdid        int64      `json:"bdid"`
}

type Board struct {
	Boardconfig string `json:"boardconfig"`
	Platform    string `json:"platform"`
	Cpid        int64  `json:"cpid"`
	Bdid        int64  `json:"bdid"`
}

type Firmware struct {
	Name		string	   `json:name`
	Identifier  string `json:"identifier"`
	Version     string     `json:"version"`
	Buildid     string     `json:"buildid"`
	Sha1Sum     string     `json:"sha1sum"`
	Md5Sum      string     `json:"md5sum"`
	Sha256Sum   string     `json:"sha256sum"`
	Filesize    int64      `json:"filesize"`
	URL         string     `json:"url"`
	Releasedate string     `json:"releasedate"`
	Uploaddate  string     `json:"uploaddate"`
	Signed      bool       `json:"signed"`
	Beta		bool	   `json:"beta"`
}

type Identifier string
const (
	IPhone81 Identifier = "iPhone8,1"
)

func decodeIpsw(data []byte) (*IPSWFirm, error) {
	var r IPSWFirm
	err := json.Unmarshal(data, &r)
	return &r, err
}

func encodeIpsw(r *IPSWFirm) ([]byte, error) {
	return json.Marshal(r)
}