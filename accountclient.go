package accountapi

import (
	"bytes"
	"fmt"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"time"
)

const (
	AccountPath = "organisation/accounts"
)

type AccountClient interface {
	Create(a *Account) (*Account, error)
	Fetch(id uuid.UUID) (*Account, error)
	List(pageNumber int, size int) ([]*Account, error)
	Delete(id uuid.UUID) error
	setTestClient(client2 *http.Client)
}

type accountClient struct {
	cl   *client
	path string
}

func NewAccountClient(endpoint *url.URL, timeout time.Duration) AccountClient {
	return &accountClient{cl: newClient(endpoint, timeout), path: AccountPath}
}

func (c *accountClient) setTestClient(t *http.Client) {
	c.cl.Client = t
}

func (c *accountClient) Create(a *Account) (*Account, error) {
	body := bytes.NewBuffer(nil)
	err := jsonapi.MarshalPayload(body, a)
	if err != nil {
		return nil, err
	}
	r, err := c.cl.post(c.path, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	a1 := new(Account)
	if err := jsonapi.UnmarshalPayload(r.Body, a1); err != nil {
		return nil, err
	}
	return a1, err
}

func (c *accountClient) Fetch(id uuid.UUID) (*Account, error) {
	r, err := c.cl.get(path.Join(c.path, id.String()))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	a := new(Account)
	if err := jsonapi.UnmarshalPayload(r.Body, a); err != nil {
		return nil, err
	}
	return a, err
}

func (c *accountClient) List(pageNumber int, size int) ([]*Account, error) {
	query := url.Values{}
	if pageNumber != 0 {
		query.Set("page[number]", fmt.Sprintf("%d", pageNumber))
	}
	if size != 0 {
		query.Set("page[size]", fmt.Sprintf("%d", size))
	}
	r, err := c.cl.get(path.Join(c.path, query.Encode()))
	if err != nil {
		return nil, err
	}

	as, err := readAccounts(r.Body)
	defer r.Body.Close()

	return as, err
}

func (c *accountClient) Delete(id uuid.UUID) error {
	r, err := c.cl.delete(path.Join(c.path, id.String()))
	if err != nil {
		return err
	}

	return r.Body.Close()
}

func readAccounts(b io.Reader) ([]*Account, error) {
	as, err := jsonapi.UnmarshalManyPayload(b, reflect.TypeOf(new(Account)))
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, len(as))
	for _, a := range as {
		if a1, ok := a.(*Account); ok {
			accounts = append(accounts, a1)
		}
	}
	return accounts, nil
}
