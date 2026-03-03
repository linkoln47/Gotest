package main

import "net/http"

type Controller struct {
	l     Logger
	logic Logic
}

func NewController(l Logger, logic Logic) Controller {
	return Controller{
		l:     l,
		logic: logic,
	}
}

func (c Controller) SayHello(w http.ResponseWriter, r *http.Request) {
	c.l.Log("In SayHello")
	userID := r.URL.Query().Get("user_id")
	message, err := c.logic.SayHello(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}

func (c Controller) SayGoodbye(w http.ResponseWriter, r *http.Request) {
	c.l.Log("In SayGoodbye")
	userID := r.URL.Query().Get("user_id")
	message, err := c.logic.SayGoodbye(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}
