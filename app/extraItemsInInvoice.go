package sdatcrm

/*
R Get("/api/extraItemsInInvoice/") -> Get all extra items in this invoice for which an adHoc order should be created
*/

import (
	"fmt"
	"net/http"

	"appengine"
)

const EXTRA_ITEMS_IN_INVOICE_API = "/api/extraItemsInInvoice/"

func init() {
	http.Handle(EXTRA_ITEMS_IN_INVOICE_API, gaeHandler(extraItemsInInvoiceHandler))
	return
}

func extraItemsInInvoiceHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pid := r.URL.Path[len(EXTRA_ITEMS_IN_INVOICE_API):]
	if len(pid) == 0 {
		switch r.Method {
		case "POST":
			invoice, err := decodeInvoice(r.Body)
			if err != nil {
				return nil, err
			}
			return FindExtraItemsInInvoice(c, invoice)
		default:
			return nil, fmt.Errorf(r.Method + " on " + r.URL.Path + " not implemented")
		}
	}
	return nil, nil
}

func FindExtraItemsInInvoice(c appengine.Context, invoice *Invoice) ([]Item, error) {
	pendingOrders, err := getPendingOrdersForPurchaser(c, invoice.PurchaserId)
	if err != nil {
		return nil, err
	}

	var clubbedPendingItems []Item
	for _, o := range pendingOrders {
		for _, i := range o.PendingItems {
			foundSameItem := false
			for _, npi := range clubbedPendingItems {
				if npi.equals(i) {
					foundSameItem = true
					npi.Qty += i.Qty
					break
				}
			}

			if !foundSameItem {
				clubbedPendingItems = append(clubbedPendingItems, i)
			}
		}
	}

	var invoicedExtraItems []Item
	for _, extraItem := range invoice.Items {
		var newItem = extraItem
		for _, i := range clubbedPendingItems {
			if extraItem.equals(i) {
				newItem.Qty -= i.Qty
				break
			}
		}
		invoicedExtraItems = append(invoicedExtraItems, newItem)
	}

	var prunedExtraItems []Item
	for _, e := range invoicedExtraItems {
		if e.Qty != 0 {
			prunedExtraItems = append(prunedExtraItems, e)
		}
		if e.Qty < 0 {
			return nil, fmt.Errorf("We should not have reached here. How can an extra item be negative. Its a bug.") // Defensive Programming.
		}
	}

	return prunedExtraItems, nil
}
