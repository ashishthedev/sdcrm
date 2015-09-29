package TODO

/*
C Post("/api/orders/") -> Create a new order
R Get("/api/orders/id/") -> Get the order with this id
U Put("/api/orders/id/") -> Resave this order with id
D Delete("/api/orders/id/") -> Delete order having id
Q Get("/api/orders/") -> Get all orders
*/

import (
	"appengine"
	"fmt"
	"net/http"
)

const ORDERS_API = "/api/orders/"

func init() {
	http.Handle(ORDERS_API, gaeHandler(orderWithID))
}

func orderWithID(c appengine.Context, w http.ResponseWriter, r *http.Request) error {
	oid := r.URL.Path[len(ORDERS_API):]
	c.Errorf("Received oid %v", oid)
	if len(oid) > 0 {
	} else {
		switch r.Method {
		case "POST":
			WriteJson(w, "Will create the order here.TODO")
		default:
			return fmt.Errorf(r.Method + " on " + r.URL.Path + " not implemented")
		}
	}
	return nil
}
