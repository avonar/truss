package test

import (
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	// 3d Party
	"golang.org/x/net/context"
	// This Service
	pb "github.com/TuneLab/go-truss/cmd/_integration-tests/transport/transport-service"
	httpclient "github.com/TuneLab/go-truss/cmd/_integration-tests/transport/transport-service/generated/client/http"

	"github.com/pkg/errors"
)

var httpAddr string

func TestGetWithQueryClient(t *testing.T) {
	var req pb.GetWithQueryRequest
	req.A = 12
	req.B = 45360
	want := req.A + req.B

	svchttp, err := httpclient.New(httpAddr)
	if err != nil {
		t.Fatalf("failed to create httpclient: %q", err)
	}

	resp, err := svchttp.GetWithQuery(context.Background(), &req)
	if err != nil {
		t.Fatalf("httpclient returned error: %q", err)
	}

	if resp.V != want {
		t.Fatalf("Expect: %d, got %d", want, resp.V)
	}
}

func TestGetWithQueryRequest(t *testing.T) {
	var resp pb.GetWithQueryResponse

	var A, B int64
	A = 12
	B = 45360
	want := A + B

	testHTTP := func(bodyBytes []byte, method, routeFormat string, routeFields ...interface{}) {
		respBytes, err := httpRequestBuilder{
			method: method,
			route:  fmt.Sprintf(routeFormat, routeFields...),
			body:   bodyBytes,
		}.Test(t)
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot make http request"))
		}

		err = json.Unmarshal(respBytes, &resp)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "json error, got response: %q", string(respBytes)))
		}

		if resp.V != want {
			t.Fatalf("Expect: %d, got %d", want, resp.V)
		}
	}

	testHTTP(nil, "GET", "getwithquery?%s=%d&%s=%d", "A", A, "B", B)
}

func TestGetWithRepeatedQueryClient(t *testing.T) {
	var req pb.GetWithRepeatedQueryRequest
	req.A = []int64{12, 45360}
	want := req.A[0] + req.A[1]

	svchttp, err := httpclient.New(httpAddr)
	if err != nil {
		t.Fatalf("failed to create httpclient: %q", err)
	}

	resp, err := svchttp.GetWithRepeatedQuery(context.Background(), &req)
	if err != nil {
		t.Fatalf("httpclient returned error: %q", err)
	}

	if resp.V != want {
		t.Fatalf("Expect: %d, got %d", want, resp.V)
	}
}

func TestGetWithRepeatedQueryRequest(t *testing.T) {
	var resp pb.GetWithRepeatedQueryResponse

	var A []int64
	A = []int64{12, 45360}
	want := A[0] + A[1]

	testHTTP := func(bodyBytes []byte, method, routeFormat string, routeFields ...interface{}) {
		respBytes, err := httpRequestBuilder{
			method: method,
			route:  fmt.Sprintf(routeFormat, routeFields...),
			body:   bodyBytes,
		}.Test(t)
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot make http request"))
		}

		err = json.Unmarshal(respBytes, &resp)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "json error, got response: %q", string(respBytes)))
		}

		if resp.V != want {
			t.Fatalf("Expect: %d, got %d", want, resp.V)
		}
	}

	testHTTP(nil, "GET", "getwithrepeatedquery?%s=[%d,%d]", "A", A[0], A[1])
	// csv style
	//testHTTP(nil, "GET", "getwithrepeatedquery?%s=%d,%d", "A", A[0], A[1])
	// multi / golang style
	//testHTTP(nil, "GET", "getwithrepeatedquery?%s=%d&%s=%d]", "A", A[0], "A", A[1])
}

func TestPostWithNestedMessageBodyClient(t *testing.T) {
	var req pb.PostWithNestedMessageBodyRequest
	var reqNM pb.NestedMessage

	reqNM.A = 12
	reqNM.B = 45360
	req.NM = &reqNM
	want := req.NM.A + req.NM.B

	svchttp, err := httpclient.New(httpAddr)
	if err != nil {
		t.Fatalf("failed to create httpclient: %q", err)
	}

	resp, err := svchttp.PostWithNestedMessageBody(context.Background(), &req)
	if err != nil {
		t.Fatalf("httpclient returned error: %q", err)
	}

	if resp.V != want {
		t.Fatalf("Expect: %d, got %d", want, resp.V)
	}
}

func TestPostWithNestedMessageBodyRequest(t *testing.T) {
	var resp pb.PostWithNestedMessageBodyResponse

	var A, B int64
	A = 12
	B = 45360
	want := A + B

	testHTTP := func(bodyBytes []byte, method, routeFormat string, routeFields ...interface{}) {
		respBytes, err := httpRequestBuilder{
			method: method,
			route:  fmt.Sprintf(routeFormat, routeFields...),
			body:   bodyBytes,
		}.Test(t)
		if err != nil {
			t.Fatal(errors.Wrap(err, "cannot make http request"))
		}

		err = json.Unmarshal(respBytes, &resp)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "json error, got response: %q", string(respBytes)))
		}

		if resp.V != want {
			t.Fatalf("Expect: %d, got %d", want, resp.V)
		}
	}

	jsonStr := fmt.Sprintf(`{ "NM": { "A": %d, "B": %d}}`, A, B)

	testHTTP([]byte(jsonStr), "POST", "postwithnestedmessagebody")
}

func TestCtxToCtxViaHTTPHeaderClient(t *testing.T) {
	var req pb.MetaRequest
	var key, value = "Truss-Auth-Header", "SECRET"
	req.Key = key

	// Create a new client telling it to send "Truss-Auth-Header" as a header
	svchttp, err := httpclient.New(httpAddr,
		httpclient.CtxValuesToSend(key))
	if err != nil {
		t.Fatalf("failed to create httpclient: %q", err)
	}

	// Create a context with the header key and value
	ctx := context.WithValue(context.Background(), key, value)

	// send the context
	resp, err := svchttp.CtxToCtx(ctx, &req)
	if err != nil {
		t.Fatalf("httpclient returned error: %q", err)
	}

	if resp.V != value {
		t.Fatalf("Expect: %q, got %q", value, resp.V)
	}
}

func TestCtxToCtxViaHTTPHeaderRequest(t *testing.T) {
	var resp pb.MetaResponse
	var key, value = "Truss-Auth-Header", "SECRET"

	jsonStr := fmt.Sprintf(`{ "Key": %q }`, key)
	fmt.Println(jsonStr)

	req, err := http.NewRequest("POST", httpAddr+"/"+"ctxtoctx", strings.NewReader(jsonStr))
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot construct http request"))
	}

	req.Header.Set(key, value)

	respBytes, err := testHTTPRequest(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot make http request"))
	}

	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "json error, got response: %q", string(respBytes)))
	}

	if resp.V != value {
		t.Fatalf("Expect: %q, got %q", value, resp.V)
	}
}

func TestErrorRPCReturnsJSONError(t *testing.T) {
	req, err := http.NewRequest("GET", httpAddr+"/"+"error", strings.NewReader(""))
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot construct http request"))
	}

	respBytes, err := testHTTPRequest(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "cannot make http request"))
	}

	jsonOut := make(map[string]interface{})
	err = json.Unmarshal(respBytes, &jsonOut)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "json error, got response: %q", string(respBytes)))
	}

	if jsonOut["error"] == nil {
		t.Fatal("http transport did not send error as json")
	}
}

// Helpers

type httpRequestBuilder struct {
	method string
	route  string
	body   []byte
}

func (h httpRequestBuilder) Test(t *testing.T) ([]byte, error) {
	t.Logf("Method: %q | Route: %q", h.method, h.route)
	httpReq, err := http.NewRequest(h.method, httpAddr+"/"+h.route, bytes.NewReader(h.body))
	if err != nil {
		return nil, err
	}

	return testHTTPRequest(httpReq)
}

func testHTTPRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not end http request")
	}
	defer httpResp.Body.Close()

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read http body")
	}

	return respBytes, nil
}
