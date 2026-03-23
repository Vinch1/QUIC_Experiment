package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type proposalRequest struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

func main() {
	var (
		addr    = flag.String("addr", "127.0.0.1:9001", "node address")
		command = flag.String("command", "put", "command name")
		key     = flag.String("key", "demo", "key")
		value   = flag.String("value", "hello", "value")
	)

	flag.Parse()

	baseURL := normalizeAddr(*addr)

	switch *command {
	case "put":
		if err := put(baseURL, *key, *value); err != nil {
			log.Fatalf("put request failed: %v", err)
		}
	case "get":
		if err := get(baseURL, *key); err != nil {
			log.Fatalf("get request failed: %v", err)
		}
	case "status":
		if err := status(baseURL); err != nil {
			log.Fatalf("status request failed: %v", err)
		}
	default:
		log.Fatalf("unsupported command: %s", *command)
	}
}

func normalizeAddr(addr string) string {
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}
	return "http://" + addr
}

func put(baseURL, key, value string) error {
	payload := proposalRequest{
		Command: "put",
		Key:     key,
		Value:   value,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/kv", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return printResponse(resp)
}

func get(baseURL, key string) error {
	resp, err := http.Get(baseURL + "/kv?key=" + key)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return printResponse(resp)
}

func status(baseURL string) error {
	resp, err := http.Get(baseURL + "/status")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return printResponse(resp)
}

func printResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("status=%d\n", resp.StatusCode)
	fmt.Print(string(body))
	return nil
}
