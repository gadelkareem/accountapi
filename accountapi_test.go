package accountapi_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gadelkareem/accountapi/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/url"
)

var f struct {
	c accountapi.AccountClient
	a *accountapi.Account
}

func form3APIServerIsRunningOn(arg1 string) error {
	u, err := url.Parse(arg1)
	if err != nil {
		return err
	}
	f.c = accountapi.NewAccountClient(u, 0)
	return nil
}

func servicesCreateAnAccountUsingData(arg1 *gherkin.DocString) error {
	f.a = new(accountapi.Account)
	err := json.Unmarshal([]byte(arg1.Content), f.a)
	if err != nil {
		return err
	}

	return nil
}

func theServerSavesAndReturnsTheAccountInTheResponse() error {
	a, err := f.c.Create(f.a)
	if err != nil {
		return err
	}
	if !assert.ObjectsAreEqual(f.a, a) {
		return errors.New("saved account does not match returned account in the response")
	}
	return nil
}

func servicesFetchAnAccountUsingAccountID(arg1 string) error {
	id, err := uuid.Parse(arg1)
	if err != nil {
		return err
	}
	f.a, err = f.c.Fetch(id)
	return err
}

func theServerReturnsTheAccountWithID(arg1 string) error {
	if f.a.ID != arg1 {
		return fmt.Errorf("account ID %s from the server does not equal account ID %s", f.a.ID, arg1)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^Form3 API server is running on "([^"]*)"$`, form3APIServerIsRunningOn)
	s.Step(`^services create an account using data$`, servicesCreateAnAccountUsingData)
	s.Step(`^the server saves and returns the account in the response$`, theServerSavesAndReturnsTheAccountInTheResponse)
	s.Step(`^services fetch an account using account ID "([^"]*)"$`, servicesFetchAnAccountUsingAccountID)
	s.Step(`^the server returns the account with ID "([^"]*)"$`, theServerReturnsTheAccountWithID)
}
