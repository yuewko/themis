package health

// TODO(miek): enable again if middleware gets health check.
/*
func TestHealth(t *testing.T) {
	h := health{Addr: ":0"}
	h.h = append(h.h, &erratic.Erratic{})

	if err := h.Startup(); err != nil {
		t.Fatalf("Unable to startup the health server: %v", err)
	}
	defer h.Shutdown()

	// Reconstruct the http address based on the port allocated by operating system.
	address := fmt.Sprintf("http://%s%s", h.ln.Addr().String(), path)

	// Norhing set should be unhealthy
	response, err := http.Get(address)
	if err != nil {
		t.Fatalf("Unable to query %s: %v", address, err)
	}
	if response.StatusCode != 503 {
		t.Errorf("Invalid status code: expecting '503', got '%d'", response.StatusCode)
	}
	response.Body.Close()

	// Make healthy
	h.Poll()

	response, err = http.Get(address)
	if err != nil {
		t.Fatalf("Unable to query %s: %v", address, err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Invalid status code: expecting '200', got '%d'", response.StatusCode)
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Unable to get response body from %s: %v", address, err)
	}
	response.Body.Close()

	if string(content) != ok {
		t.Errorf("Invalid response body: expecting 'OK', got '%s'", string(content))
	}
}
*/
