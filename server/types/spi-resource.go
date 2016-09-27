package types

// Mfrc522Resource : mfrc522Resource resource struct
type Mfrc522Sector16 struct {
	Values [16]byte `json:"values"`
}

// Mfrc522Resource : mfrc522Resource resource struct
type Mfrc522Resource struct {
	Key    [6]byte `json:"key"`
	Uid    [5]byte `json:"uid"`
	Status int     `json:"status"`
	Len    int     `json:"len"`
	Data   []byte  `json:"data"`
	// DumpClassic1K request
	Sectors [64]Mfrc522Sector16 `json:"sectors"`
}
