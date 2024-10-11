package main

import (
	"log"
	"net/http"
	"time"

	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav/carddav"
)

func main() {
	addr := "0.0.0.0:8888"
	backend := bwbBackend{
		currentUserPrincipal:   "/ZGU/",
		addressBookHomeSetPath: "contacts",
		addresses: []carddav.AddressObject{
			{
				Path:    "/ZGU/contacts/default/alice.vcf",
				ModTime: time.Now(),
				ETag:    "sQNI/mCtS7HUYkK+447YWozmRX10-Fest",
				Card: map[string][]*vcard.Field{
					"FN":      {{Group: "", Value: "Fest T"}},
					"N":       {{Group: "", Value: "T;Fest;;;"}},
					"PRODID":  {{Group: "", Value: "-//Apple Inc.//iOS 17.5.1//EN"}},
					"REV":     {{Group: "", Value: "2024-08-25T08:27:42Z"}},
					"UID":     {{Group: "", Value: "A4DCAEA8-996C-4113-AD57-4D66BC05E986"}},
					"VERSION": {{Group: "", Value: "3.0"}},
				},
			},
		},
	}
	handler := carddav.Handler{
		Backend: &backend,
		Prefix:  "",
	}

	log.Printf("WebDAV server listening on %v", addr)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		handler.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServeTLS(addr, "localhost.crt", "localhost.key", h))
}
