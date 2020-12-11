package models

//PaymentUpdateResponse gives paymentupdateresponse of records
type PaymentUpdateResponse struct {
	Records *StatusRecord `json:"Records"`
}

//StatusRecord gives status of record
type StatusRecord struct {
	Status string `json:"Status"`
}
