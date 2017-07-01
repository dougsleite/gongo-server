package main

import (
	"net/http"
	"gopkg.in/mgo.v2"
)

func main() {
	session := mongoSession("localhost")
	defer session.Close()

	router := NewRouter(session)

	router.Handle("query execute", queryExecute)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}

func mongoSession(url string) (*mgo.Session) {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}
