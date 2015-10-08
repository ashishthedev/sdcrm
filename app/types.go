package sdatcrm

import (
	"time"
)

type TCDieItem struct {
	PelletSize string `json:"PelletSize"`
	BoreSize   string `json:"BoreSize"`
	CaseType   string `json:"CaseType"`
	CaseSize   string `json:"CaseSize"`
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
	SKU
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

type Order struct {
	Id              OrderId     `json:"Id" datastore:"-"`
	Created         time.Time   `json:"Created"`
	Date            time.Time   `json:"Date"`
	TotalQty        int64       `json:"TotalQty"`
	PurchaserId     PurchaserId `json:"PurchaserId"`
	SupplierName    string      `json:"SupplierName"`
	Number          string      `json:"Number"`
	Pending         bool        `json:"Pending"`
	InvoicesId      []InvoiceId `json:"InvoicesId"`
	OrderedItems    []Item      `json:"OrderedItems"`
	PendingItems    []Item      `json:"PendingItems"`
	DispatchedItems []Item      `json:"DispatchedItems"`
	PuntedItems     []Item      `json:"PuntedItems"`
}

type Invoice struct {
	Id             InvoiceId `json:"Id" datastore:"-"`
	Items          []Item    `json:"Items"`
	Created        time.Time `json:"Created"`
	Date           time.Time `json:"Date"`
	TotalQty       int64     `json:"TotalQty"`
	PurchaserName  string    `json:"PurchaserName"`
	SupplierName   string    `json:"SupplierName"`
	Number         string    `json:"Number"`
	PRemarks       string    `json:"PRemarks"`
	OrdersId       []OrderId `json:"OrdersId"`
	DoNotMoveStock bool      `json:"DoNotMoveStock"`
}
