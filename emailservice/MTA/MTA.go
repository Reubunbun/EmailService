package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"strings"
	"io/ioutil"
	"io"
	"os"
	"time"
	"bytes"
	"log"
)

type Email struct {
	From string
	To   string
	Body string
}

var serverAddress string

func SendEmail() {
	index    := -1
	maxIndex := 1
	client   := &http.Client {}
	var email Email
	
	for {
	
		time.Sleep( 1000 * time.Millisecond )
		if index >= maxIndex - 1 {
			index = 0
		} else {
			index += 1
		}
		
		// List the users on this email server
		usersUrl := "http://" + serverAddress + "1:8888/users"
		if users, getUsersErr := MakeRequest( client, "GET", usersUrl, nil ); getUsersErr == nil && users != "" {
			userList    := strings.Split( users, "," )
			maxIndex     = len( userList ) - 1
			currentUser := userList[ index ]
			titlesUrl   := "http://" + serverAddress + "1:8888/users/" + currentUser + "/outbox"
			
			// List the emails in this users outbox
			if titles, getTitlesErr := MakeRequest( client, "GET", titlesUrl, nil ); getTitlesErr == nil && titles != "[]" {
				firstTitle := strings.Split( titles, "," )[0]
				// Remove array brackets
				firstTitle  = strings.Replace( firstTitle, "[", "", -1 )
				firstTitle  = strings.Replace( firstTitle, "]", "", -1 )
				// Remove string quotes
				firstTitle  = strings.Replace( firstTitle, "\"", "", -1 )
				emailUrl   := "http://" + serverAddress + "2:8888/MSA/" + currentUser + "/outbox/" + firstTitle
				
				// Get the first email listed
				if emailString, getEmailErr := MakeRequest( client, "GET", emailUrl, nil ); getEmailErr == nil {
					emailSplit  := strings.Split( emailString, "\n" )
					to          := strings.Split( emailSplit[0], "To:   " )[1]
					from        := strings.Split( emailSplit[1], "From: " )[1]
					body        := strings.Split( emailSplit[2], "Body: " )[1]
					email        = Email{ To : to, From : from, Body : body }
					bluebookUrl := "http://192.168.1.2:8888/bluebook/" + strings.Split( to, "@" )[1]
					
					// Find the ip address of the destination email server
					if address, getAddressErr := MakeRequest( client, "GET", bluebookUrl, nil ); getAddressErr == nil {
						MTAUrl := "http://" + address + "3:8888/MTA/recieveEmail/" + firstTitle
							
						// POST this email to the destination servers MTA
						if enc, encErr := json.Marshal( email ); encErr == nil {
							if _, putErr := MakeRequest( client, "PUT", MTAUrl, bytes.NewBuffer( enc ) ); putErr == nil {
								fmt.Printf( "Email sent\n" )
								deleteUrl  := "http://" + serverAddress + "2:8888/MSA/" + currentUser + "/outbox/" + firstTitle
					
								// Delete the email
								if _, deleteErr := MakeRequest( client, "DELETE", deleteUrl, nil ); deleteErr != nil {
									fmt.Println( "DELETE failed with %s\n", deleteErr )
								}	
									
							} else {
								fmt.Printf( "POST failed with %s\n", putErr )
								continue
							}
								
						} else {
							fmt.Printf( "Encode failed with %s\n", encErr )
							continue
						}
						
					} else {
						fmt.Println( "Address GET failed with %s\n", getAddressErr )
						continue
					}				
					
				} else {
					fmt.Printf( "Email GET failed with %s\n", getEmailErr )
					continue
				}
				
			} else if getTitlesErr != nil {
				fmt.Printf( "Titles GET failed with %s\n", getTitlesErr )
				continue
			} else if titles == "[]" {
				continue
			}
			
		} else if getUsersErr != nil {
			fmt.Printf( "User GET failed with %s\n", getUsersErr )
			continue
		} else if users == "" {
			continue
		}
		
	}
}

func RecieveEmail( w http.ResponseWriter, r *http.Request ) {
	time.Sleep( 2000 * time.Millisecond )
	client  := &http.Client {}
	vars    := mux.Vars( r )
	subject := vars[ "subject" ]
	url     := "http://" + serverAddress + "2:8888/MSA/inbox/" + subject
	decoder := json.NewDecoder( r.Body )
	var email Email
	
	if decodeErr := decoder.Decode( &email ); decodeErr == nil {
		fmt.Println( email )
		if enc, encErr := json.Marshal( email ); encErr == nil {
			if resp, postErr := MakeRequest( client, "PUT", url, bytes.NewBuffer( enc ) ); postErr == nil {
				fmt.Fprintln( w, resp )
			} else {
				fmt.Printf( "PUT failed with %s\n", postErr )
			}
		} else {
			fmt.Printf( "Encode failed with %s\n", encErr )
		}
	} else {
		fmt.Printf( "Decode failed with %s\n", decodeErr )
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
	router.HandleFunc( "/MTA/recieveEmail/{subject}", RecieveEmail ).Methods( "PUT" )
	log.Fatal( http.ListenAndServe( ":8888", router ) )
}

func main(){
	serverAddress = os.Args[1:][0]
	go SendEmail()
	time.Sleep(time.Millisecond)
	HandleRequests()
}



