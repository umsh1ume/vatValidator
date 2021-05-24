package services

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	ViesURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService.com"

	checkVatRequestFormat = `
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
	<soapenv:Header/>
	<soapenv:Body>
	  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
	    <countryCode><<.countryCode>></countryCode>
	    <vatNumber><<.vatNumber>></vatNumber>
	  </checkVat>
	</soapenv:Body>
	</soapenv:Envelope>
	`
	GermanVatFormatRegex = `^(DE){1}[0-9]{9}$`

	InvalidGermanFormatErrorMsg = `A German VAT number starts with DE followed by 9 numeric characters.`

	InvalidVatOnViesErrorMsg = `The VAT number is invalid on VIES`
)

type CheckVatResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    struct {
		XMLName xml.Name `xml:"Body"`
		Soap    struct {
			XMLName     xml.Name `xml:"checkVatResponse"`
			CountryCode string   `xml:"countryCode"`
			VATNumber   string   `xml:"vatNumber"`
			RequestDate string   `xml:"requestDate"` // 2015-03-06+01:00
			Valid       bool     `xml:"valid"`
			Name        string   `xml:"name"`
			Address     string   `xml:"address"`
		}
	}
}

func ValidateGermanVat(vatNumber string) (bool, string, error) {

	isValid, msg, err := validateGermanFormat(vatNumber)

	if err != nil || !isValid {
		return isValid, msg, err
	}

	// If valid german format, then check it on Vies
	isValid, msg, err = validateVatOnVies(vatNumber)

	return isValid, msg, err
}

func validateGermanFormat(vatNumber string) (bool, string, error) {
	matched, err := regexp.Match(GermanVatFormatRegex, []byte(vatNumber))

	if err != nil {
		return false, "", err
	}

	// If the vat number does not match german format, return false with custom message
	if !matched {
		return false, InvalidGermanFormatErrorMsg, nil
	}

	return true, "", nil
}

func validateVatOnVies(vatNumber string) (bool, string, error) {
	viesReq := getViesRequestPayload(vatNumber)
	viesReqBytes := bytes.NewBufferString(viesReq)
	client := http.Client{}

	res, err := client.Post(ViesURL, "text/xml;charset=UTF-8", viesReqBytes)

	if err != nil {
		return false, "", err
	}
	defer res.Body.Close()

	xmlRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, "", err
	}

	var viesResponse CheckVatResponse

	if err = xml.Unmarshal(xmlRes, &viesResponse); err != nil {
		return false, "", err
	}

	if viesResponse.Soap.Soap.Valid != true {
		return false, InvalidVatOnViesErrorMsg, nil
	}

	return true, "", nil
}

func getViesRequestPayload(n string) string {
	countryCode := n[0:2]
	vatNumber := n[2:]

	req := checkVatRequestFormat
	req = strings.Replace(req, "<<.countryCode>>", countryCode, 1)
	req = strings.Replace(req, "<<.vatNumber>>", vatNumber, 1)
	return req
}
