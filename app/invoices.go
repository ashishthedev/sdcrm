package sdatcrm

/*
C Post("/api/invoices/") -> Create a new invoice
R Get("/api/invoices/id/") -> Get the invoice with this id
U Put("/api/invoices/id/") -> Resave this invoice with id
D Delete("/api/invoices/id/") -> Delete invoice having id
Q Get("/api/invoices/") -> Get all invoices
*/

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
)

const INVOICES_API = "/api/invoices/"

type InvoiceId int64

func init() {
	http.Handle(INVOICES_API, gaeHandler(invoiceHandler))
	http.HandleFunc("/invoice/new/", newInvoicePageHandler)
	http.HandleFunc("/invoice/", editInvoicePageHandler)
	http.HandleFunc("/invoices/", allInvoicePageHandler)
}

func invoiceHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id := r.URL.Path[len(INVOICES_API):]
	c.Infof("Received invoice id %v", id)
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
		default:
			return nil, fmt.Errorf(r.Method + " on " + r.URL.Path + " not implemented")
		}
	} else {
		switch r.Method {
		case "POST":
			return invoiceSaveEntryPoint(c, r)
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

func newInvoicePageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/invoice.html"))
	var data interface{}
	data = struct{ Nature string }{"NEW"}
	if t == nil {
		t = PAGE_NOT_FOUND_TEMPLATE
		data = nil
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func editInvoicePageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/invoice.html"))
	var data interface{}
	data = struct{ Nature string }{"EDIT"}
	if t == nil {
		t = PAGE_NOT_FOUND_TEMPLATE
		data = nil
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func allInvoicePageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/invoices.html"))
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func UpdateRelatedOrdersFieldAndEditTheirPendingItemsListAndMarkOrderCompletion(c appengine.Context, invoice *Invoice) error {

	return nil
}

func invoiceSaveEntryPoint(c appengine.Context, r *http.Request) (*Invoice, error) {
	invoice, err := decodeInvoice(r.Body)
	if err != nil {
		return nil, err
	}
	if invoice.Id == 0 {
		return ProcessBrandNewInvoice(c, invoice)
	} else {
		return ProcessBrandNewInvoice(c, invoice)
		//return ProcessExistingInvoice(c, invoice)
	}
}

func ProcessBrandNewInvoice(c appengine.Context, invoice *Invoice) (*Invoice, error) {
	var newInvoice *Invoice
	var err error
	err = datastore.RunInTransaction(c, func(c appengine.Context) error {
		if err = CreateAdHocOrderIfRequired(c, invoice); err != nil {
			return err
		}

		if err = UpdateRelatedOrdersFieldAndEditTheirPendingItemsListAndMarkOrderCompletion(c, invoice); err != nil {
			return err
		}

		newInvoice, err = invoice.save(c)
		return err
	}, nil)

	if err != nil {
		return nil, err
	} else {
		return newInvoice, nil
	}
}

func ProcessExistingInvoice(c appengine.Context, invoice *Invoice) (*Invoice, error) {
	return invoice.save(c)
}

func CreateAdHocOrderWithTheseItems(c appengine.Context, items []Item, invoice *Invoice) (*Order, error) {
	order := new(Order)
	order.PurchaserId = invoice.PurchaserId
	order.SupplierName = invoice.SupplierName
	order.Date = time.Now()
	order.Number = "Telephonic"
	for _, i := range items {
		order.OrderedItems = append(order.OrderedItems, i)
		order.PendingItems = append(order.PendingItems, i)
	}
	c.Infof("About to save order:%#v", order)
	return order.save(c)
}
func CreateAdHocOrderIfRequired(c appengine.Context, invoice *Invoice) error {
	// 1. Check if the invoice is being created for some extra items.
	// 2. Create an adHoc order for extra items.
	// 3. Recheck if invoice is being created for extra items. This time it should not be. Defensive Programming.
	extraItems, err := FindExtraItemsInInvoice(c, invoice)
	if err != nil {
		return err
	}
	c.Infof("Extra items in invoice of %v: %#v", invoice.PurchaserId, extraItems)

	if len(extraItems) > 0 {
		o, err := CreateAdHocOrderWithTheseItems(c, extraItems, invoice)
		if err != nil {
			return err
		}
		c.Infof("created teh extra order: %#v", o)

	}
	return nil
}
