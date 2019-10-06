package accountapi

type Account struct {
	ID                          string   `jsonapi:"primary,accounts"`
	Country                     string   `jsonapi:"attr,country"`
	BaseCurrency                string   `jsonapi:"attr,base_currency"`
	AccountNumber               string   `jsonapi:"attr,account_number"`
	BankID                      string   `jsonapi:"attr,bank_id"`
	BankIDCode                  string   `jsonapi:"attr,bank_id_code"`
	Bic                         string   `jsonapi:"attr,bic"`
	Iban                        string   `jsonapi:"attr,iban"`
	CustomerId                  string   `jsonapi:"attr,customer_id"`
	Title                       string   `jsonapi:"attr,title"`
	FirstName                   string   `jsonapi:"attr,first_name"`
	BankAccountName             string   `jsonapi:"attr,bank_account_name"`
	AlternativeBankAccountNames []string `jsonapi:"attr,alternative_bank_account_names"`
	AccountClassification       string   `jsonapi:"attr,account_classification"`
	JointAccount                bool     `jsonapi:"attr,joint_account"`
	AccountMatchingOptOut       bool     `jsonapi:"attr,account_matching_opt_out"`
	SecondaryIdentification     string   `jsonapi:"attr,secondary_identification"`
}
