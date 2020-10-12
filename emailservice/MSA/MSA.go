package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil" 
	"io"
	"os"
	"log"
	"fmt"
	"bytes"
	"strings"
)

type Email struct {
	From string
	To   string
	Body string
}

var serverAddress string

func ToBox( w http.ResponseWriter, r *http.Request ) {
	client  := &http.Client {}
	vars    := mux.Vars( r )
	box     := vars[ "box"     ]
	subject := vars[ "subject" ]
	decoder := json.NewDecoder( r.Body )
	var email Email
	if decodeErr := decoder.Decode( &email ); decodeErr == nil {
		
		var url string
		if box == "inbox" {
			url = "http://" + serverAddress + "1:8888/users/" + email.To + "/inbox/" + subject
		} else if box == "outbox" {
			url = "http://" + serverAddress + "1:8888/users/" + email.From + "/outbox/" + subject
		} else {
			w.WriteHeader( http.StatusBadRequest )
		}
		
		if enc, encErr := json.Marshal( email ); encErr == nil {
			if resp, putErr := MakeRequest( client, "PUT", url, bytes.NewBuffer( enc ) ); putErr == nil {
				fmt.Fprintln( w, resp )
			} else {
				fmt.Fprintf( w, "PUT failed with %s\n", putErr )
			}
		} else {
			w.WriteHeader( http.StatusInternalServerError )
		}
		
		
	} else {
		w.WriteHeader( http.StatusBadRequest )
	}
}


func ListBox( w http.ResponseWriter, r *http.Request ) {
	client := &http.Client {}
	vars   := mux.Vars( r )
	user   := vars[ "user" ]
	box    := vars[ "box"  ]
	url    := "http://" + serverAddress + "1:8888/users/" + user + "/" + box
	
	if body, getErr := MakeRequest( client, "GET", url, nil ); getErr == nil {
		// Replace "," with "|" to make it easier to split the list correctly (there may be commas in the title)
		titlesString := strings.Replace( body, "\",\"", "\"|\"", -1 )
		// Remove array brackets
		titlesString  = strings.Replace( titlesString, "[", "", -1 )
		titlesString  = strings.Replace( titlesString, "]", "", -1 )
		titlesList   := strings.Split( titlesString, "|" )
		
		for index, title := range titlesList{
			fmt.Fprintf( w, "(%d) " + title + "\n", index + 1 )
		}
	} else {
		fmt.Fprintf( w, "GET failed %s\n", getErr )
	}
}

func ReadEmail( w http.ResponseWriter, r *http.Request ) {
	client := &http.Client {}
	vars   := mux.Vars( r )
	user   := vars[ "user"  ]
	box    := vars[ "box"   ]
	title  := vars[ "title" ]
	url    := "http://" + serverAddress + "1:8888/users/" + user + "/" + box + "/" + title
	if body, getErr := MakeRequest( client, "GET", url, nil ); getErr == nil {
		decoder := json.NewDecoder( bytes.NewBufferString( body ) )
		var email Email
		if decodeErr := decoder.Decode( &email ); decodeErr == nil {
			fmt.Fprintf( w, "To:   %s\nFrom: %s\nBody: %s\n", email.To, email.From, email.Body )
		} else {
			w.WriteHeader( http.StatusInternalServerError )
		}
	} else {
		fmt.Fprintf( w, "GET failed with %s\n", getErr )
	}
}

func DeleteEmail( w http.ResponseWriter, r *http.Request ) {
	client := &http.Client {}
	vars   := mux.Vars( r )
	user   := vars[ "user"  ]
	box    := vars[ "box"   ]
	title  := vars[ "title" ]
	url    := "http://" + serverAddress + "1:8888/users/" + user + "/" + box + "/" + title
	
	if resp, deleteErr := MakeRequest( client, "DELETE", url, nil ); deleteErr == nil {
		fmt.Fprintln( w, resp )
	} else {
		fmt.Fprintf( w, "DELETE failed with %s\n", deleteErr )
	}
}

func MakeRequest( client *http.Client, request string, url string, data io.Reader ) ( string, error ) {
	if req, badReqErr := http.NewRequest( request, url, data ); badReqErr == nil {
		if resp, reqErr := client.Do( req ); reqErr == nil {
			if body, respErr := ioutil.ReadAll( resp.Body ); respErr == nil {
				return string( body ), respErr
			} else {
				return "", respErr
			}
		} else {
			return "", reqErr
		}
	} else {
		return "", badReqErr
	}
}


func HandleRequests() {
	router := mux.NewRouter().StrictSlash( true )
	router.HandleFunc( "/MSA/{box}/{subject}"     , ToBox       ).Methods( "PUT"    )
	router.HandleFunc( "/MSA/{user}/{box}"        , ListBox     ).Methods( "GET"    )
	router.HandleFunc( "/MSA/{user}/{box}/{title}", ReadEmail   ).Methods( "GET"    )
	router.HandleFunc( "/MSA/{user}/{box}/{title}", DeleteEmail ).Methods( "DELETE" )
	log.Fatal( http.ListenAndServe( ":8888", router ) )
}

func main() {
	serverAddress = os.Args[1:][0]
	HandleRequests()
}