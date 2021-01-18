package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var (
	Application = "http://localhost:7070"

	TestDataDir = "testdata"
	WordList    = "wordlist.json"
)

type Response struct {
	Error   string   `json:"error"`
	Success bool     `json:"success"`
	Result  []string `json:"result"`
}

func getAnagrams(word string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/get?word=%s", Application, word))
	if err != nil {
		return nil, fmt.Errorf("http.Get error: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}
	respData := Response{}
	if err = json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	if err = resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("close response body error: %v", err)
	}

	return respData.Result, nil
}

func loadWordlist(fileName string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get work dir: %w", err)
	}
	jsonFile, err := os.Open(filepath.Join(workDir, "..", TestDataDir, fileName))
	if err != nil {
		return fmt.Errorf("failed to open wordlist file: %w", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read wordlist file: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/load", Application), "application/json", bytes.NewBuffer(byteValue))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %w", err)
	}
	respData := Response{}
	if err = json.Unmarshal(body, &respData); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	if respData.Success != true {
		return fmt.Errorf("failed to load dictionary: %+v", respData)
	}
	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("close response body error: %v", err)
	}

	return nil
}

func TestGet(t *testing.T) {
	if err := loadWordlist(WordList); err != nil {
		t.Errorf("Load wordlist error: %v", err)
	}

	tests := []struct {
		in   string
		want Response
	}{
		{
			in:   fmt.Sprintf("%s/get?word=", Application),
			want: Response{Error: "couldn't found 'word' parameter", Success: false},
		},
		{
			in:   fmt.Sprintf("%s/get?TEST=TEST", Application),
			want: Response{Error: "couldn't found 'word' parameter", Success: false},
		},
		{
			in:   fmt.Sprintf("%s/get?word=DaD", Application),
			want: Response{Result: []string{"add"}, Success: true},
		},
		{
			in:   fmt.Sprintf("%s/get?word=cat", Application),
			want: Response{Result: []string{"cat", "act"}, Success: true},
		},
		{
			in:   fmt.Sprintf("%s/get?word=UnknownWordForTest", Application),
			want: Response{Result: nil, Success: true},
		},
	}

	for _, tt := range tests {
		resp, err := http.Get(tt.in)
		if err != nil {
			t.Errorf("http.Get error: %v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("read response body error: %v", err)
		}
		got := Response{}
		if err = json.Unmarshal(body, &got); err != nil {
			t.Errorf("unmarshal error: %v", err)
		}
		if err = resp.Body.Close(); err != nil {
			t.Errorf("close response body error: %v", err)
		}
		require.Equal(t, got, tt.want)
	}
}

func TestLoadWordlist(t *testing.T) {
	tests := []struct {
		inReqData    []byte
		wantResponse Response

		inWord       string
		wantAnagrams []string
	}{
		{
			inReqData:    []byte(`["invalid", "Json"`),
			wantResponse: Response{Error: "failed to decode input json: unexpected EOF", Success: false},
			inWord:       "random",
			wantAnagrams: nil,
		},
		{
			inReqData:    []byte(`["foobar", "aabb", "baba", "boofar", "test"]`),
			wantResponse: Response{Success: true},
			inWord:       "foobar",
			wantAnagrams: []string{"boofar", "foobar"},
		},
		{
			inReqData:    []byte(`["TEST", "eStt", "tset"]`),
			wantResponse: Response{Success: true},
			inWord:       "ttse",
			wantAnagrams: []string{"tset", "eStt", "TEST"},
		},
		{
			inReqData:    []byte(`["test"]`),
			wantResponse: Response{Success: true},
			inWord:       "foobar",
			wantAnagrams: nil,
		},
	}

	for _, tt := range tests {
		appHost := fmt.Sprintf("%s/load", Application)
		resp, err := http.Post(appHost, "application/json", bytes.NewReader(tt.inReqData))
		if err != nil {
			t.Errorf("failed to send POST request: %v", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("read response body error: %v", err)
		}
		gotResponse := Response{}
		if err = json.Unmarshal(body, &gotResponse); err != nil {
			t.Errorf("unmarshal error: %v", err)
		}
		if err = resp.Body.Close(); err != nil {
			t.Errorf("close response body error: %v", err)
		}
		require.Equal(t, gotResponse, tt.wantResponse)

		gotAnagrams, err := getAnagrams(tt.inWord)
		if err != nil {
			t.Errorf("getAnagrams error: %v", err)
		}
		require.Equal(t, gotAnagrams, tt.wantAnagrams)
	}
}

func TestAddWords(t *testing.T) {
	if err := loadWordlist(WordList); err != nil {
		t.Errorf("Load wordlist error: %v", err)
	}
	client := &http.Client{}

	tests := []struct {
		inReqData    []byte
		wantResponse Response

		inWord       string
		wantAnagrams []string
	}{
		{
			inReqData:    []byte(`["invalid", "Json", "tac"`),
			wantResponse: Response{Error: "failed to decode input json: unexpected EOF", Success: false},
			inWord:       "tac",
			wantAnagrams: []string{"cat", "act"},
		},
		{
			inReqData:    []byte(`[""]`),
			wantResponse: Response{Success: true},
			inWord:       "tac",
			wantAnagrams: []string{"cat", "act"},
		},
		{
			inReqData:    []byte(`["tac", "ATC"]`),
			wantResponse: Response{Success: true},
			inWord:       "tac",
			wantAnagrams: []string{"ATC", "tac", "cat", "act"},
		},
	}
	for _, tt := range tests {
		appHost := fmt.Sprintf("%s/load", Application)
		req, err := http.NewRequestWithContext(context.Background(), "PATCH", appHost, bytes.NewReader(tt.inReqData))
		if err != nil {
			t.Errorf("failed to create PATCH request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("failed to send PATCH request: %v", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("read response body error: %v", err)
		}
		gotResponse := Response{}
		if err = json.Unmarshal(body, &gotResponse); err != nil {
			t.Errorf("unmarshal error: %v", err)
		}
		if err = resp.Body.Close(); err != nil {
			t.Errorf("close response body error: %v", err)
		}

		require.Equal(t, gotResponse, tt.wantResponse)

		gotAnagrams, err := getAnagrams(tt.inWord)
		if err != nil {
			t.Errorf("getAnagrams error: %v", err)
		}
		require.Equal(t, gotAnagrams, tt.wantAnagrams)
	}
}
