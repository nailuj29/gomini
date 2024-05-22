package client

// TitanRequest sends a request to a Titan server
//
// If token or mime is not desired, an empty string can be passed
func TitanRequest(address string, body []byte, token string, mime string) Response {
	if mime == "" {
		mime = "text/gemini"
	}
}
