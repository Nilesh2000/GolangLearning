package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestSetHandler(t *testing.T) {
	store := NewKVStore()
	server := NewServer(store)

	tests := []struct {
		name         string
		method       string
		contentType  string
		body         interface{} // can be a map or string
		wantStatus   int
		wantResponse map[string]string // nil means expect empty body
	}{
		{
			name:        "Wrong HTTP Method (GET)",
			method:      http.MethodGet,
			contentType: "application/json",
			body:        map[string]string{"key": "username", "value": "gopher"},
			wantStatus:  http.StatusMethodNotAllowed,
			wantResponse: map[string]string{
				"error": "Method not allowed",
			},
		},
		{
			name:        "Wrong Content-Type",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        map[string]string{"key": "username", "value": "gopher"},
			wantStatus:  http.StatusUnsupportedMediaType,
			wantResponse: map[string]string{
				"error": "Content-Type must be application/json",
			},
		},
		{
			name:        "Invalid body",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        "{invalid}",
			wantStatus:  http.StatusBadRequest,
			wantResponse: map[string]string{
				"error": "Bad Request",
			},
		},
		{
			name:        "Missing value",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        map[string]string{"key": "username"},
			wantStatus:  http.StatusBadRequest,
			wantResponse: map[string]string{
				"error": "Key and Value are required",
			},
		},
		{
			name:        "Empty Key",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        map[string]string{"key": "", "value": "gopher"},
			wantStatus:  http.StatusBadRequest,
			wantResponse: map[string]string{
				"error": "Key and Value are required",
			},
		},
		{
			name:         "Valid Request",
			method:       http.MethodPost,
			contentType:  "application/json",
			body:         map[string]string{"key": "username", "value": "gopher"},
			wantStatus:   http.StatusCreated,
			wantResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte

			switch b := tt.body.(type) {
			case string:
				bodyBytes = []byte(b)
			case map[string]string:
				var err error
				bodyBytes, err = json.Marshal(b)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(tt.method, "/set", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			server.SetHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantResponse == nil {
				if rr.Body.Len() != 0 {
					t.Errorf("expected empty response, got %q", rr.Body.String())
				}
			} else {
				var got map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(got, tt.wantResponse) {
					t.Errorf("got %q, want %q", got, tt.wantResponse)
				}
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	store := NewKVStore()
	server := NewServer(store)

	store.Set("foo", "bar")

	tests := []struct {
		name         string
		method       string
		keyParam     string
		wantStatus   int
		wantResponse map[string]string
	}{
		{
			name:       "Wrong HTTP Method (POST)",
			method:     http.MethodPost,
			keyParam:   "username",
			wantStatus: http.StatusMethodNotAllowed,
			wantResponse: map[string]string{
				"error": "Method not allowed",
			},
		},
		{
			name:       "Missing key",
			method:     http.MethodGet,
			keyParam:   "",
			wantStatus: http.StatusBadRequest,
			wantResponse: map[string]string{
				"error": "Key is required",
			},
		},
		{
			name:       "Non-existent key",
			method:     http.MethodGet,
			keyParam:   "bar",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]string{
				"error": "Key not found",
			},
		},
		{
			name:       "Valid Request",
			method:     http.MethodGet,
			keyParam:   "foo",
			wantStatus: http.StatusOK,
			wantResponse: map[string]string{
				"key":   "foo",
				"value": "bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/get?key="+tt.keyParam, nil)
			rr := httptest.NewRecorder()

			server.GetHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d want %d", rr.Code, tt.wantStatus)
			}

			var got map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tt.wantResponse) {
				t.Errorf("got %q, want %q", got, tt.wantResponse)
			}
		})
	}
}

func TestDeleteEndpoint(t *testing.T) {
	store := NewKVStore()
	server := NewServer(store)

	store.Set("foo", "bar")

	tests := []struct {
		name         string
		method       string
		keyParam     string
		wantStatus   int
		wantResponse map[string]string
		wantDeleted  bool
	}{
		{
			name:       "Wrong HTTP Method (GET)",
			method:     http.MethodGet,
			keyParam:   "foo",
			wantStatus: http.StatusMethodNotAllowed,
			wantResponse: map[string]string{
				"error": "Method not allowed",
			},
			wantDeleted: false,
		},
		{
			name:       "Missing Key",
			method:     http.MethodDelete,
			keyParam:   "",
			wantStatus: http.StatusBadRequest,
			wantResponse: map[string]string{
				"error": "Key is required",
			},
			wantDeleted: false,
		},
		{
			name:       "Non-existent Key",
			method:     http.MethodDelete,
			keyParam:   "bar",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]string{
				"error": "Key not found",
			},
			wantDeleted: false,
		},
		{
			name:         "Valid Request",
			method:       http.MethodDelete,
			keyParam:     "foo",
			wantStatus:   http.StatusOK,
			wantResponse: nil,
			wantDeleted:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/delete?key="+tt.keyParam, nil)
			rr := httptest.NewRecorder()

			server.DeleteHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantResponse == nil {
				if rr.Body.Len() != 0 {
					t.Errorf("expected empty response, got %q", rr.Body.String())
				}
			} else {
				var got map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tt.wantResponse) {
					t.Errorf("got %q, want %q", got, tt.wantResponse)
				}
			}

			if tt.wantDeleted {
				if _, exists := store.Get(tt.keyParam); exists {
					t.Errorf("key %q should have been deleted", tt.keyParam)
				}
			}
		})
	}
}

func TestKeysEndpoint(t *testing.T) {
	store := NewKVStore()
	server := NewServer(store)

	tests := []struct {
		name         string
		method       string
		setup        func()
		wantStatus   int
		wantResponse interface{}
	}{
		{
			name:       "Wrong HTTP Method (POST)",
			method:     http.MethodPost,
			setup:      func() {},
			wantStatus: http.StatusMethodNotAllowed,
			wantResponse: map[string]string{
				"error": "Method not allowed",
			},
		},
		{
			name:       "Empty store",
			method:     http.MethodGet,
			setup:      func() {},
			wantStatus: http.StatusOK,
			wantResponse: map[string][]string{
				"keys": {},
			},
		},
		{
			name:   "Valid Request",
			method: http.MethodGet,
			setup: func() {
				store.Set("foo", "bar")
				store.Set("key", "value")
				store.Set("username", "gopher")
			},
			wantStatus: http.StatusOK,
			wantResponse: map[string][]string{
				"keys": {"foo", "key", "username"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest(tt.method, "/keys", nil)
			rr := httptest.NewRecorder()

			server.KeysHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d want %d", rr.Code, tt.wantStatus)
			}

			switch want := tt.wantResponse.(type) {
			case map[string]string:
				var got map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
				}

			case map[string][]string:
				var got map[string][]string
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				sort.Strings(got["keys"])
				sort.Strings(want["keys"])
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		})
	}
}
