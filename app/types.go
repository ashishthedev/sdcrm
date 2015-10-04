package sdatcrm

type TCDieItem struct {
	PelletSize string  `json:"PelletSize"`
	BoreSize   float64 `json:"BoreSize"`
	CaseType   string  `json:"CaseType"`
	CaseSize   string  `json:"CaseSize"`
}

type MiscItem struct {
	Name string `json:"Name"`
	Unit string `json:"Unit"`
}

type SKU struct {
	TCDieItem
	MiscItem
	Rate     float64 `json:"Rate"`
	Type     string  `json:"Type"` //TCD or MSC
	CRemarks string  `json:"CRemarks"`
}

type Item struct {
	SKU SKU
	Qty int64 `json:"Qty"`
}
