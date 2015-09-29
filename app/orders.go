package sdatcrm

/*
C Post("/api/orders/") -> Create a new order
R Get("/api/orders/id/") -> Get the order with this id
U Put("/api/orders/id/") -> Resave this order with id
D Delete("/api/orders/id/") -> Delete order having id
Q Get("/api/orders/") -> Get all orders
*/

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ORDERS_API = "/api/orders/"

type Item struct {
	Type         string  `json:"Type"`
	FriendlyName string  `json:"FriendlyName"`
	Remarks      string  `json:"Remarks"`
	PelletSize   string  `json:"PelletSize"`
	Unit         string  `json:"Unit"`
	BoreSize     float64 `json:"BoreSize"`
	CasingType   string  `json:"CasingType"`
	CasingSize   string  `json:"CasingSize"`
	Qty          int64   `json:"Qty"`
}

type Order struct {
	Id            int64     `json:"Id" datastore:"-"`
	Items         []Item    `json:"Items"`
	Created       time.Time `json:"Created"`
	PODate        time.Time `json:"PODate"`
	TotalQty      int64     `json:"TotalQty"`
	PurchaserName string    `json:"PurchaserName"`
	SupplierName  string    `json:"SupplierName"`
	PONumber      string    `json:"PONumber"`
	Done          bool      `json:"Done"`
}

func init() {
	http.Handle(ORDERS_API, gaeHandler(orderWithID))
}

func orderWithID(c appengine.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	oid := r.URL.Path[len(ORDERS_API):]
	c.Errorf("Received oid %v", oid)
	if len(oid) > 0 {
	} else {
		switch r.Method {
		case "POST":
			order, err := decodeOrder(r.Body)
			if err != nil {
				return nil, err
			}
			return order.save(c)
		default:
			return nil, fmt.Errorf(r.Method + " on " + r.URL.Path + " not implemented")
		}
	}
	return nil, nil
}

func decodeOrder(r io.ReadCloser) (*Order, error) {
	defer r.Close()
	var order Order
	err := json.NewDecoder(r).Decode(&order)
	return &order, err
}

func (o *Order) save(c appengine.Context) (*Order, error) {
	k, err := datastore.Put(c, o.key(c), o)
	if err != nil {
		return nil, err
	}
	o.Id = k.IntID()
	return o, nil
}

func defaultOrdersList(c appengine.Context) *datastore.Key {
	ancestorKey := datastore.NewKey(c, "ANCESTOR_KEY", BranchName(c), 0, nil)
	return datastore.NewKey(c, "OrderList", "default", 0, ancestorKey)
}

func (o *Order) key(c appengine.Context) *datastore.Key {
	if o.Id == 0 {
		o.Created = time.Now()
		return datastore.NewIncompleteKey(c, "Order", defaultOrdersList(c))
	}
	return datastore.NewKey(c, "Order", "", o.Id, defaultOrdersList(c))
}
