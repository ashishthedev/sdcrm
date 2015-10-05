package sdatcrm

import (
	"time"
)

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

type Address struct {
	DeliveryAddress      string `json:"DeliveryAddress"`
	City                 string `json:"City"`
	State                string `json:"State"`
	Pincode              string `json:"Pincode"`
	EnvelopePhoneNumbers CSL    `json:"EnvelopePhoneNumbers"`
}

type Purchaser struct {
	Address
	Id              PurchaserId `json:"Id" datastore:"-"`
	SKUs            []SKU       `json:"SKUs"`
	Created         time.Time   `json:"Created"`
	Name            string      `json:"Name"`
	DispatchEmails  CSL         `json:"DispatchEmails"`
	FORMCEmails     CSL         `json:"FORMCEmails"`
	TinNumber       string      `json:"TinNumber"`
	BillingAddress  string      `json:"BillingAddress"`
	SMSPhoneNumbers CSL         `json:"SMSPhoneNumbers"`
	MDPhoneNumbers  CSL         `json:"MDPhoneNumbers"`
	CreditDays      int64       `json:"CreditDays"`
	CRemarks        string      `json:"CRemarks"`
}
