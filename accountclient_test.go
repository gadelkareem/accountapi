package accountapi

import (
	"bytes"
	"context"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"path/filepath"
	"testing"
)

const jsonPath = "data/json"

func TestAccountClient_Fetch(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	b := readFile(t, "accounts/fetch01.json")
	testClient, teardown := testingHTTPAccountClient(
		t,
		id.String(),
		http.MethodGet,
		http.StatusOK,
		b,
	)
	defer teardown()

	c := newAccountTestClient(testClient)

	a1, err := c.Fetch(id)
	checkError(t, err)
	assert.Nil(t, err)

	a2 := new(Account)
	err = jsonapi.UnmarshalPayload(bytes.NewBuffer(b), a2)
	checkError(t, err)
	assert.Equal(t, a1, a2)
}

func BenchmarkAccountClient_Fetch(b *testing.B) {
	id := uuid.New()
	testClient, teardown := benchmarkingHTTPAccountClient(http.StatusOK, readFile(b, "accounts/fetch01.json"))
	defer teardown()
	c := newAccountTestClient(testClient)
	for n := 0; n < b.N; n++ {
		_, err := c.Fetch(id)
		checkError(b, err)
	}
}

func TestAccountClient_List(t *testing.T) {
	t.Parallel()
	b := readFile(t, "accounts/list01.json")
	filters := url.Values{}
	filters.Set("page[number]", "1")
	filters.Set("page[size]", "100")
	testClient, teardown := testingHTTPAccountClient(
		t,
		filters.Encode(),
		http.MethodGet,
		http.StatusOK,
		b,
	)
	defer teardown()

	c := newAccountTestClient(testClient)

	as1, err := c.List(1, 100)
	checkError(t, err)
	assert.Nil(t, err)

	as2, err := readAccounts(bytes.NewBuffer(b))
	checkError(t, err)
	assert.Equal(t, as1, as2)
}

func BenchmarkAccountClient_List(b *testing.B) {
	testClient, teardown := benchmarkingHTTPAccountClient(http.StatusOK, readFile(b, "accounts/list01.json"))
	defer teardown()
	c := newAccountTestClient(testClient)
	for n := 0; n < b.N; n++ {
		_, err := c.List(0, 0)
		checkError(b, err)
	}
}

func TestAccountClient_Create(t *testing.T) {
	t.Parallel()
	b := readFile(t, "accounts/create01.json")
	testClient, teardown := testingHTTPAccountClient(
		t,
		"",
		http.MethodPost,
		http.StatusCreated,
		b,
	)
	defer teardown()

	c := newAccountTestClient(testClient)

	a2 := new(Account)
	err := jsonapi.UnmarshalPayload(bytes.NewBuffer(b), a2)
	checkError(t, err)

	a1, err := c.Create(a2)
	checkError(t, err)
	assert.Nil(t, err)

	assert.Equal(t, a1, a2)
}

func BenchmarkAccountClient_Create(b *testing.B) {
	body := readFile(b, "accounts/create01.json")
	testClient, teardown := benchmarkingHTTPAccountClient(http.StatusOK, body)
	defer teardown()
	a := new(Account)
	err := jsonapi.UnmarshalPayload(bytes.NewBuffer(body), a)
	checkError(b, err)
	c := newAccountTestClient(testClient)
	for n := 0; n < b.N; n++ {
		_, err := c.Create(a)
		checkError(b, err)
	}
}

func TestAccountClient_Delete(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	testClient, teardown := testingHTTPAccountClient(
		t,
		id.String(),
		http.MethodDelete,
		http.StatusNoContent,
		nil,
	)
	defer teardown()

	c := newAccountTestClient(testClient)

	err := c.Delete(id)
	checkError(t, err)
}

func BenchmarkAccountClient_Delete(b *testing.B) {
	id := uuid.New()
	testClient, teardown := benchmarkingHTTPAccountClient(http.StatusNoContent, nil)
	defer teardown()
	c := newAccountTestClient(testClient)
	for n := 0; n < b.N; n++ {
		err := c.Delete(id)
		checkError(b, err)
	}
}

func readFile(t assert.TestingT, f string) []byte {
	pth, err := filepath.Abs(filepath.Join(jsonPath, f))
	checkError(t, err)
	b, err := ioutil.ReadFile(pth)
	checkError(t, err)

	return b
}

func checkError(t assert.TestingT, err error) {
	if err != nil {
		t.Errorf("%s", err)
	}
}

func newAccountTestClient(t *http.Client) AccountClient {
	endpoint, _ := url.Parse("http://test/")
	c := NewAccountClient(endpoint, 0)
	c.setTestClient(t)
	return c
}

func benchmarkingHTTPAccountClient(statusCode int, body []byte) (*http.Client, func()) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if body != nil {
			w.Write(body)
		}
	}))

	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return c, s.Close
}

func testingHTTPAccountClient(t assert.TestingT, urlPath, method string, statusCode int, body []byte) (*http.Client, func()) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/vnd.api+json", r.Header.Get("Accept"))
		assert.Equal(t, r.URL.Path, path.Join("/", Version, AccountPath, urlPath))
		assert.Equal(t, r.Method, method)
		w.WriteHeader(statusCode)
		if body != nil {
			w.Write(body)
		}
	}))

	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return c, s.Close
}
