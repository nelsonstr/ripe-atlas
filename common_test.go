package atlas

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetVersion(t *testing.T) {
	ver := GetVersion()
	assert.EqualValues(t, ourVersion, ver, "should be equal")
}

func TestGetPageNum(t *testing.T) {
	url := "https://foo.example.com/"

	n := getPageNum(url)
	if n != "" {
		t.Errorf("n=%s should be ''", n)
	}
	url = "https://foo.example.com/?page=42"
	n = getPageNum(url)
	if n != "42" {
		t.Errorf("n=%s should be 42", n)
	}
	url = "https://foo.example.com/?country=fr&page=43"
	n = getPageNum(url)
	if n != "43" {
		t.Errorf("n=%s should be 43", n)
	}
	url = "https://foo.example.com/?country=fr&page=666&bar=1"
	n = getPageNum(url)
	if n != "666" {
		t.Errorf("n=%s should be 666", n)
	}

	url = ""
	n = getPageNum(url)
	if n != "" {
		t.Errorf("n=%s should be ''", n)
	}
}

func TestClienthandleAPIResponsese(t *testing.T) {
	var (
		r http.Response
		b bytes.Buffer
	)

	client, err := NewClient()
	err = client.handleAPIResponsese(nil)
	assert.Error(t, err, "should be in error")

	r = http.Response{StatusCode: 0}
	err = client.handleAPIResponsese(&r)
	assert.NoError(t, err, "should be no error")

	r = http.Response{StatusCode: 200}
	err = client.handleAPIResponsese(&r)
	assert.NoError(t, err, "should be no error")

	var jsonErr = `error:{status: 501, code: 500, detail: "test"}`

	fmt.Fprintf(&b, "%v", jsonErr)
	r.StatusCode = 300
	r.Body = ioutil.NopCloser(&b)
	err = client.handleAPIResponsese(&r)
	assert.NoError(t, err, "should be in error")

	r.StatusCode = 500
	r.Body = ioutil.NopCloser(&b)
	err = client.handleAPIResponsese(&r)
	assert.Error(t, err, "should be in error")
}

func TestAddQueryParameters(t *testing.T) {
	p := AddQueryParameters("", map[string]string{})
	assert.Equal(t, "", p)
}

func TestAddQueryParameters_1(t *testing.T) {
	p := AddQueryParameters("", map[string]string{"": ""})
	assert.Equal(t, "?=", p)
}

func TestAddQueryParameters_2(t *testing.T) {
	p := AddQueryParameters("foo", map[string]string{"bar": "baz"})
	assert.Equal(t, "foo?bar=baz", p)
}

func TestClient_AddAPIKey(t *testing.T) {
	c, err := NewClient(Config{APIKey:"foo"})
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotEmpty(t, c)

	opts := map[string]string{}

	new := c.addAPIKey(opts)
	assert.NotEmpty(t, c.config.APIKey)
	assert.Equal(t, 1, len(new))
	assert.EqualValues(t, map[string]string{"key": "foo"}, new)
}
