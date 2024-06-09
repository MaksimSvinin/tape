package model

type TapeInfo struct {
	Vendor   string `json:"vendor,omitempty"`
	Model    string `json:"model,omitempty"`
	Firmware string `json:"firmware,omitempty"`

	Attributes *Attributes `json:"attributes,omitempty"`
}
