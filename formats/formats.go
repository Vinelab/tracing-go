package formats

var (
	// TextMap is a format descriptor propagating trace context via a map
	TextMap = "text_map"

	// HTTP is a format descriptor for propagating trace context via HTTP request
	HTTP = "http"

	// AMQP is a format descriptor for propagating trace context via AMQP message
	AMQP = "amqp"
)
