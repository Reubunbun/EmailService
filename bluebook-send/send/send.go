package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"io"
	"strings"
	"log"
	"fmt"
	"bytes"
)

type Email struct {
	From string
	To   string
	Body string
}

func Send( w http.ResponseWriter, r *http.Request ) {
	client  := &http.Client {}
	decoder := json.NewDecoder( r.Body )
	vars    := mux.Vars( r )
	subject := vars[ "subject" ]
	var email Email
	if err := decoder.Decode( &email ); err == nil {
		server  := strings.Split( email.From, "@" )[1]
		url     := "http://192.168.1.2:8888/bluebook/" + server
		if address, getErr := MakeRequest( client, "GET", url, nil ); getErr == nil {
			url = "http://" + address + "2:8888/MSA/outbox/" + subject
			if enc, encErr := json.Marshal( email ); encErr == nil {
				if resp, postErr := MakeRequest( client, "PUT", url, bytes.NewBuffer( enc ) ); postErr == nil {
					fmt.Fprintln( w, resp )
				} else {
					fmt.Fprintf( w, "PUT failed with %s\n", postErr )
				}
			} else {
				w.WriteHeader( http.StatusInternalServerError )
			}
		} else {
			fmt.Fprintf( w, "GET failed with %s\n", getErr )
		}
	} else {
		w.WriteHeader( http.StatusBadRequest )
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
	router.HandleFunc( "/send/{subject}", Send ).Methods( "PUT" )
	log.Fatal( http.ListenAndServe( ":8888", router ) )
}

func main() {
	HandleRequests()
}