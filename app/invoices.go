package sdatcrm

/*
C Post("/api/invoices/") -> Create a new invoice
R Get("/api/invoices/id/") -> Get the invoice with this id
U Put("/api/invoices/id/") -> Resave this invoice with id
D Delete("/api/invoices/id/") -> Delete invoice having id
Q Get("/api/invoices/") -> Get all invoices
*/

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const INVOICES_API = "/api/invoices/"

type Invoice struct {
	Id            InvoiceId `json:"Id" datastore:"-"`
	Items         []Item    `json:"Items"`
	Created       time.Time `json:"Created"`
	Date          time.Time `json:"Date"`
	TotalQty      int64     `json:"TotalQty"`
	PurchaserName string    `json:"PurchaserName"`
	SupplierName  string    `json:"SupplierName"`
	Number        string    `json:"Number"`
	OrdersId      []OrderId `json:"OrdersId"`
}

func init() {
	http.Handle(INVOICES_API, gaeHandler(invoiceHandler))
}

func invoiceHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := r.URL.Path[len(INVOICES_API):]
	c.Errorf("Received invoice id %v", id)
	if len(id) > 0 {
		switch r.Method {
		case "GET":
			id64, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, err
			}
			invoice := new(Invoice)
			invoice.Id = InvoiceId(id64)
			return invoice.get(c)
		}
	} else {
		switch r.Method {
		case "POST":
			invoice, err := decodeInvoice(r.Body)
			if err != nil {
				return nil, err
			}
			return invoice.save(c)
		case "GET":
			return getAllInvoices(c)
		default:
			return nil, fmt.Errorf(r.Method + " on " + r.URL.Path + " not implemented")
		}
	}
	return nil, nil
}

func decodeInvoice(r io.ReadCloser) (*Invoice, error) {
	defer r.Close()
	var invoice Invoice
	err := json.NewDecoder(r).Decode(&invoice)
	return &invoice, err
}

func (o *Invoice) get(c appengine.Context) (*Invoice, error) {
	err := datastore.Get(c, o.key(c), o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
func (o *Invoice) save(c appengine.Context) (*Invoice, error) {
	k, err := datastore.Put(c, o.key(c), o)
	if err != nil {
		return nil, err
	}
	o.Id = InvoiceId(k.IntID())
	return o, nil
}

func defaultInvoiceList(c appengine.Context) *datastore.Key {
	ancestorKey := datastore.NewKey(c, "ANCESTOR_KEY", BranchName(c), 0, nil)
	return datastore.NewKey(c, "InvoiceList", "default", 0, ancestorKey)
}

func (o *Invoice) key(c appengine.Context) *datastore.Key {
	if o.Id == 0 {
		o.Created = time.Now()
		return datastore.NewIncompleteKey(c, "Invoice", defaultInvoiceList(c))
	}
	return datastore.NewKey(c, "Invoice", "", int64(o.Id), defaultInvoiceList(c))
}

func getAllInvoices(c appengine.Context) ([]Invoice, error) {
	invoices := []Invoice{}
	ks, err := datastore.NewQuery("Invoice").Ancestor(defaultInvoiceList(c)).Order("Created").GetAll(c, &invoices)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(invoices); i++ {
		invoices[i].Id = InvoiceId(ks[i].IntID())
	}
	return invoices, nil
}
