package types

// DioResource : dio resource struct
type DioResource struct {
	Pin        int    `json:"pin"`
	Sender     uint64 `json:"sender"`
	Interuptor uint64 `json:"interruptor"`
	On         bool   `json:"on"`
}
