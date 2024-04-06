package esi

// esiStatusRequest is a request to the ESI for a status.
func esiStatusRequest() request {
	return newRequest("/latest/status/", make(map[string]string), 0)
}
