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

// Mfrc522DumpResource : Mfrc522DumpResource resource struct
type Mfrc522DumpResource struct {
	Key     [6]byte `json:"key"`
	Uid     [5]byte `json:"uid"`
	TagType string  `json:"tagType"`
	// DumpClassic1K request
	Sectors []Mfrc522Sector16 `json:"sectors"`
}

// Mfrc522WriteResource : Mfrc522WriteResource resource struct
type Mfrc522WriteResource struct {
	Key     [6]byte `json:"key"`
	Uid     [5]byte `json:"uid"`
	TagType string  `json:"tagType"`
	// WriteClassic1K request
	Sector byte            `json:"sector"`
	Data   Mfrc522Sector16 `json:"sector"`
}
