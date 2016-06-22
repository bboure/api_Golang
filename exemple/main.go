package main

import (
	"log"

	"github.com/PlanetHoster/api_Golang/phapi"
)

const (
	APIKEY  = "mykey"
	APIUSER = "myuser"
)

func main() {
	sld := "exemple"
	tld := "com"
	//Aucune validation n'est faite à ce point
	api := phapi.New(APIKEY, APIUSER)

	//Vérifie si la connexion est valide
	if err := api.Test(); err != nil {
		log.Fatalln("Erreur de connexion", err)
	} else {
		log.Println("Connexion réussi")
	}

	//Retrouve les informations de votre compte revendeur
	log.Println(api.AccountInfo())

	//Vérifie la disponibiltié d'un domain
	log.Println(api.DomainAvailable(sld, tld))

	//Plus d'information sur un domaine
	log.Println(api.DomainInfo(sld, tld))

	//Retourne les informations du WHOIS du domaine
	log.Println(api.Whois(sld, tld))

	//Retourne les nameservers du domaine
	log.Println(api.Nameservers(sld, tld))

	//Si les DNS sont avec PlanetHoster, les retournes
	log.Println(api.DNSRecords(sld, tld))

	//Listes toutes les extensions disponibles
	log.Println(api.TLDPrices())

	//Enregistre un domaine (commenté pour éviter les erreurs)

	/*	domain := phapi.NewDomainData(&phapi.ContactDomain{
			FirstName:   "John",
			LastName:    "Bob",
			Email:       "test@monmail.com",
			CompanyName: "",
			Address1:    "123 Chemin des Planète",
			Address2:    "",
			City:        "Paris",
			PostalCode:  "1234567",
			State:       "Paris",
			CountryCode: "FR",
			Phone:       "0176604143",
		}, "nsa.planethoster.net")
		//Il est recommandé d'utiliser un minimum de 2 NS
		domain.NS2 = "nsb.planethoster.net"
		domain.NS3 = "nsc.planethoster.net"

		fmt.Println(api.RegisterDomain(sld, tld, 1, domain))
	*/
}
