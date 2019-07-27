package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"./core"
)

func Test_getJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(fullChain))
	defer ts.Close()

	type args struct {
		url    string
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", args{ts.URL, &core.Blockchain{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getJSON(tt.args.url, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("getJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerNodes(t *testing.T) {
	//創建一個請求
	data := url.Values{}
	data.Add("addr", "http://localhost:6060")
	data.Add("addr", "http://localhost:6070")

	req, err := http.NewRequest("POST", "/nodes/register", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 我們創建一個 ResponseRecorder (which satisfies http.ResponseWriter)來記錄響應
	rr := httptest.NewRecorder()

	//直接使用HealthCheckHandler，傳入參數 rr,req
	registerNodes(rr, req)

	// 檢測返回的狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢測返回的數據
	expected := "http://localhost:6060,\nhttp://localhost:6070,\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_fullChain(t *testing.T) {
	//創建一個請求
	req, err := http.NewRequest("GET", "/chain/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 我們創建一個 ResponseRecorder (which satisfies http.ResponseWriter)來記錄響應
	rr := httptest.NewRecorder()

	//直接使用HealthCheckHandler，傳入參數 rr,req
	fullChain(rr, req)

	// 檢測返回的狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢測返回的數據
	expected := `{"chain":[{"index":1,"timestamp"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_newTransactions(t *testing.T) {
	//創建一個請求
	data := url.Values{}
	data.Set("sender", "zps")
	data.Set("recipient", "sun")
	data.Set("amount", "123.2")

	req, err := http.NewRequest("POST", "/transactions/new", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 我們創建一個 ResponseRecorder (which satisfies http.ResponseWriter)來記錄響應
	rr := httptest.NewRecorder()

	//直接使用HealthCheckHandler，傳入參數 rr,req
	newTransactions(rr, req)

	// 檢測返回的狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢測返回的數據
	expected := `[{s:zps, r:sun, $123.20}]`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_mine(t *testing.T) {
	//創建一個請求
	req, err := http.NewRequest("GET", "/mine", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 我們創建一個 ResponseRecorder (which satisfies http.ResponseWriter)來記錄響應
	rr := httptest.NewRecorder()

	//直接使用HealthCheckHandler，傳入參數 rr,req
	mine(rr, req)

	// 檢測返回的狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢測返回的數據
	expected := "[{\n#1"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_consensus(t *testing.T) {
	// create a listener with the desired port.
	l, err := net.Listen("tcp", "127.0.0.1:6060")
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(http.HandlerFunc(fullChain))

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener.Close()
	ts.Listener = l

	// Start the server.
	ts.Start()

	// Run your tests here.

	// Cleanup.
	defer ts.Close()

	//創建一個請求
	req, err := http.NewRequest("GET", "/nodes/resolve", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 我們創建一個 ResponseRecorder (which satisfies http.ResponseWriter)來記錄響應
	rr := httptest.NewRecorder()

	//直接使用HealthCheckHandler，傳入參數 rr,req
	consensus(rr, req)

	// 檢測返回的狀態碼
	if status := rr.Code; status != http.StatusOK && status != http.StatusRequestTimeout {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢測返回的數據
	expected := "[{\n#1"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
