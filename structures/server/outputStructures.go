package serverStructures

// ServerRejection
/*
 * Caller is the function that invokes the request.
 * Content is the description of the request, or what triggered the rejection.
 * Rejection is the error itself, to be accessed by the handler function.
 * ServerRejection passes the meaningful content to be handled by a helper function.
 * It is used for functions or methods that fail (request parsing, message sending, et cetera)
 */
type ServerRejection struct {
	Caller    string
	Content   string
	Rejection error
}

// ServerOutput
/*
 * Caller is the function that invokes the request.
 * Content is the information to be transmitted.
 * ServerOutput passes the meaningful content to be handled by a helper function.
 * It is used for regular output (request logging, latency checks, et cetera)
 */
type ServerOutput struct {
	Caller  string
	Content string
}

// ServerWarning
/*
 * Caller is the function that invokes the request.
 * Content is the information to be transmitted.
 * ServerWarning passes the meaningful content to be handled by a helper function.
 * It is used for soft errors, which don't generate exceptions (insufficient arguments, for instance)
 */
type ServerWarning struct {
	Caller  string
	Content string
}
