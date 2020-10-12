package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"os"
	"net/http"
	"fmt"
	"time"
)

type Email struct {
	From string
	To   string
	Body string
}

var serverName string
var users map[ string ] map[ string ] map[ string ] Email

func Create( w http.ResponseWriter, r *http.Request ) {
	if body, err := ioutil.ReadAll( r.Body ); err == nil {
		user := string( body ) + "@" + serverName
		w.WriteHeader( http.StatusCreated )
		users[ user ] = make( map[ string ] map[ string ] Email )
		users[ user ][ "inbox"  ] = make( map[ string ] Email )
		users[ user ][ "outbox" ] = make( map[ string ] Email )
	} else {
		w.WriteHeader( http.StatusBadRequest )
	}
} 

func Store( w http.ResponseWriter, r *http.Request ) {
	vars    := mux.Vars( r )
	user    := vars[ "user"    ]
	box     := vars[ "box"     ]
	subject := vars[ "subject" ]
	decoder := json.NewDecoder( r.Body )
	var email Email
	if err := decoder.Decode( &email ); err == nil {
		if _, userExists := users[ user ][ box ]; userExists {
			if box == "outbox" {
				dt           := time.Now()
				dateFormat   := dt.Format("01:02:2006-15:04:05")
				subjWithDate := subject + "-" + dateFormat
				
				if _, emailExists := users[ user ][ box ][ subjWithDate ]; emailExists {
					dateFormat   = dt.Format("01-02-2006-15:04:05.000000")
					subjWithDate = subject + "-" + dateFormat
					
					if _, emailStillExists := users[ user ][ box ][ subjWithDate ]; emailStillExists {
						w.WriteHeader( http.StatusInternalServerError )
					} else {
						w.WriteHeader( http.StatusCreated  )
						users[ user ][ box ][ subjWithDate ] = email
					}
				} else {
					w.WriteHeader( http.StatusCreated  )
					users[ user ][ box ][ subjWithDate ] = email
				}
			} else if box == "inbox" {
				w.WriteHeader( http.StatusCreated  )
				users[ user ][ box ][ subject ] = email
			} else {
				w.WriteHeader( http.StatusBadRequest )
			}
		} else {
			w.WriteHeader( http.StatusNotFound )
		}
	} else {
		w.WriteHeader( http.StatusBadRequest )
	}
}

func List( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	user := vars[ "user" ]
	box  := vars[ "box"  ]
	if emails, exist := users[ user ][ box ]; exist {
		titles := []string{}
		for email, _ := range emails {
			titles = append( titles, email )
		}
		if enc, err := json.Marshal( titles ); err == nil {
			w.WriteHeader( http.StatusOK )
			w.Write( []byte( enc ) )
		} else {
			w.WriteHeader( http.StatusInternalServerError )
		}
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}

func Read( w http.ResponseWriter, r *http.Request ) {
	vars  := mux.Vars( r )
	user  := vars[ "user"  ]
	box   := vars[ "box"   ]
	title := vars[ "title" ]
	if email, exists := users[ user ][ box ][ title ]; exists {
		if enc, err := json.Marshal( email ); err == nil {
			w.WriteHeader( http.StatusOK )
			w.Write( []byte( enc ) )
		} else {
			w.WriteHeader( http.StatusInternalServerError )
		}
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
	
}

func ReadUsers( w http.ResponseWriter, r *http.Request ) {
	w.WriteHeader( http.StatusOK )
	usersList := ""
	for user := range users {
		usersList = usersList + user + ","
	}
	fmt.Fprintf( w, usersList )
}

func Delete( w http.ResponseWriter, r *http.Request ) {
	vars  := mux.Vars( r )
	user  := vars[ "user"  ]
	box   := vars[ "box"   ]
	title := vars[ "title" ]
	
	if emails, userExists := users[ user ][ box ]; userExists {
		if _, emailExists := emails[ title ]; emailExists {
			w.WriteHeader( http.StatusNoContent )
			delete( emails, title )
			users[ user ][ box ] = emails
		} else {
			w.WriteHeader( http.StatusBadRequest )
		}
	} else {
		w.WriteHeader( http.StatusBadRequest )
	}
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash( true )
	router.HandleFunc( "/users"                       , Create    ).Methods( "POST"    )
	router.HandleFunc( "/users/{user}/{box}/{subject}", Store     ).Methods( "PUT"    )
	router.HandleFunc( "/users/{user}/{box}/{title}"  , Delete    ).Methods( "DELETE" )
	router.HandleFunc( "/users/{user}/{box}/{title}"  , Read      ).Methods( "GET"    )
	router.HandleFunc( "/users/{user}/{box}"          , List      ).Methods( "GET"    )
	router.HandleFunc( "/users"                       , ReadUsers ).Methods( "GET"    )
	log.Fatal( http.ListenAndServe( ":8888", router ) )	
}

func main() {
	serverName = os.Args[1:][0]
	users = make( map[ string ] map[ string ] map[ string ] Email )
	HandleRequests()
}








