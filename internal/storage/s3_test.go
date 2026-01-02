package storage

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type roundTripper func(*http.Request) (*http.Response, error)

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func newTestS3Storage(t *testing.T, handler func(*http.Request) (*http.Response, error)) *S3Storage {
	t.Helper()

	cfg := aws.Config{
		Region: "us-east-1",
		HTTPClient: &http.Client{
			Transport: roundTripper(handler),
		},
		Credentials: aws.AnonymousCredentials{},
	}

	return &S3Storage{
		client: s3.NewFromConfig(cfg),
		bucket: "test-bucket",
	}
}

func TestS3Storage_Save(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	err := s.Save("k", []byte("v"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestS3Storage_Save_Error(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("put error")
	})

	err := s.Save("k", []byte("v"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestS3Storage_Load(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("data"))),
		}, nil
	})

	data, err := s.Load("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "data" {
		t.Fatalf("unexpected data")
	}
}

func TestS3Storage_Load_Error(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("get error")
	})

	_, err := s.Load("k")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestS3Storage_Delete(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		return &http.Response{
			StatusCode: 204,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	err := s.Delete("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestS3Storage_Delete_Error(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("delete error")
	})

	err := s.Delete("k")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestS3Storage_Exists_True(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodHead {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	ok, err := s.Exists("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected exists")
	}
}

func TestS3Storage_Exists_False(t *testing.T) {
	s := newTestS3Storage(t, func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("not found")
	})

	ok, err := s.Exists("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected not exists")
	}
}
