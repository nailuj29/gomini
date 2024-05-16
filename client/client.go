package client

// Response represents a Gemini Response
type Response struct {
	// Data contains the raw data returned from the server.
	// TODO: implement function to pipe this straight to Gemtext parser
	Data []byte
	// ContentType contains the MIME type of the Data
	ContentType string
}

func Request(url string) Response {
	return Response{}
}
