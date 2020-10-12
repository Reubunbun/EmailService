package main

import  ( 
	"github.com/gorilla/mux"
	"io/ioutil" 
	"log" 
	"net/http" 
)

var servers map[ string ] string

func Create( w http.ResponseWriter, r *http.Request ) {
	vars   := mux.Vars( r )
	server := vars[ "server" ]
	if body, err := ioutil.ReadAll( r.Body ); err == nil {
		w.WriteHeader( http.StatusCreated )
		address := string( body )
		servers[ server ] = address
	} else {
		w.WriteHeader( http.StatusBadRequest )
	}
}

func Read( w http.ResponseWriter, r *http.Request ) {
	vars   := mux.Vars( r )
	server := vars[ "server" ]
	if address, ok := servers[ server ]; ok {
		w.WriteHeader( http.StatusOK )
		w.Write( []byte( address ) )
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}

func Update( w http.ResponseWriter, r *http.Request ) { 
	vars := mux.Vars( r ) 
	server := vars[ "server" ] 
	if _, ok := servers[ server ]; ok { 
		if body, err := ioutil.ReadAll( r.Body ); err == nil {
			w.WriteHeader( http.StatusCreated ) 
			address := string( body )
			servers[ server ] = address
		} else {
			w.WriteHeader( http.StatusBadRequest )
		}		
	} else { 
		w.WriteHeader( http.StatusNotFound ) 
	} 
}

func Delete( w http.ResponseWriter, r *http.Request ) {
	vars := mux.Vars( r )
	server := vars[ "server" ]
	if _, ok := servers[ server ]; ok {
		w.WriteHeader( http.StatusNoContent )
		delete( servers, server )
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash( true )
	router.HandleFunc( "/bluebook/{server}", Create ).Methods( "POST"   )
	router.HandleFunc( "/bluebook/{server}", Read   ).Methods( "GET"    )
	router.HandleFunc( "/bluebook/{server}", Update ).Methods( "PUT"    )
	router.HandleFunc( "/bluebook/{server}", Delete ).Methods( "DELETE" )
	log.Fatal( http.ListenAndServe( ":8888", router ) )
}

func main() {
	servers = make( map[ string ] string )
	HandleRequests()
}