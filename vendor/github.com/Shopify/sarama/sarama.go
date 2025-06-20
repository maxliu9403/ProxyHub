/*
Package sarama is a pure Go client library for dealing with Apache Kafka (versions 0.8 and later). It includes a high-level
API for easily producing and consuming messages, and a low-level API for controlling bytes on the wire when the high-level
API is insufficient. Usage examples for the high-level APIs are provided inline with their full documentation.

To produce messages, use either the AsyncProducer or the SyncProducer. The AsyncProducer accepts messages on a channel
and produces them asynchronously in the background as efficiently as possible; it is preferred in most cases.
The SyncProducer provides a method which will block until Kafka acknowledges the message as produced. This can be
useful but comes with two caveats: it will generally be less efficient, and the actual durability guarantees
depend on the configured value of `Producer.RequiredAcks`. There are configurations where a message acknowledged by the
SyncProducer can still sometimes be lost.

To consume messages, use Consumer or Consumer-Group API.

For lower-level needs, the Broker and Request/Response objects permit precise control over each connection
and message sent on the wire; the Client provides higher-level metadata management that is shared between
the producers and the consumer. The Request/Response objects and properties are mostly undocumented, as they line up
exactly with the protocol fields documented by Kafka at
https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol

Metrics are exposed through https://github.com/rcrowley/go-metrics library in a local registry.

Broker related metrics:

	+----------------------------------------------+------------+---------------------------------------------------------------+
	| Name                                         | Type       | Description                                                   |
	+----------------------------------------------+------------+---------------------------------------------------------------+
	| incoming-byte-rate                           | meter      | Bytes/second read off all brokers                             |
	| incoming-byte-rate-for-broker-<broker-id>    | meter      | Bytes/second read off a given broker                          |
	| outgoing-byte-rate                           | meter      | Bytes/second written off all brokers                          |
	| outgoing-byte-rate-for-broker-<broker-id>    | meter      | Bytes/second written off a given broker                       |
	| request-rate                                 | meter      | Requests/second sent to all brokers                           |
	| request-rate-for-broker-<broker-id>          | meter      | Requests/second sent to a given broker                        |
	| request-size                                 | histogram  | Distribution of the request size in bytes for all brokers     |
	| request-size-for-broker-<broker-id>          | histogram  | Distribution of the request size in bytes for a given broker  |
	| request-latency-in-ms                        | histogram  | Distribution of the request latency in ms for all brokers     |
	| request-latency-in-ms-for-broker-<broker-id> | histogram  | Distribution of the request latency in ms for a given broker  |
	| response-rate                                | meter      | Responses/second received from all brokers                    |
	| response-rate-for-broker-<broker-id>         | meter      | Responses/second received from a given broker                 |
	| response-size                                | histogram  | Distribution of the response size in bytes for all brokers    |
	| response-size-for-broker-<broker-id>         | histogram  | Distribution of the response size in bytes for a given broker |
	| requests-in-flight                           | counter    | The current number of in-flight requests awaiting a response  |
	|                                              |            | for all brokers                                               |
	| requests-in-flight-for-broker-<broker-id>    | counter    | The current number of in-flight requests awaiting a response  |
	|                                              |            | for a given broker                                            |
	+----------------------------------------------+------------+---------------------------------------------------------------+

Note that we do not gather specific metrics for seed brokers but they are part of the "all brokers" metrics.

Producer related metrics:

	+-------------------------------------------+------------+--------------------------------------------------------------------------------------+
	| Name                                      | Type       | Description                                                                          |
	+-------------------------------------------+------------+--------------------------------------------------------------------------------------+
	| batch-size                                | histogram  | Distribution of the number of bytes sent per partition per request for all topics    |
	| batch-size-for-topic-<topic>              | histogram  | Distribution of the number of bytes sent per partition per request for a given topic |
	| record-send-rate                          | meter      | Records/second sent to all topics                                                    |
	| record-send-rate-for-topic-<topic>        | meter      | Records/second sent to a given topic                                                 |
	| records-per-request                       | histogram  | Distribution of the number of records sent per request for all topics                |
	| records-per-request-for-topic-<topic>     | histogram  | Distribution of the number of records sent per request for a given topic             |
	| compression-ratio                         | histogram  | Distribution of the compression ratio times 100 of record batches for all topics     |
	| compression-ratio-for-topic-<topic>       | histogram  | Distribution of the compression ratio times 100 of record batches for a given topic  |
	+-------------------------------------------+------------+--------------------------------------------------------------------------------------+

Consumer related metrics:

	+-------------------------------------------+------------+--------------------------------------------------------------------------------------+
	| Name                                      | Type       | Description                                                                          |
	+-------------------------------------------+------------+--------------------------------------------------------------------------------------+
	| consumer-batch-size                       | histogram  | Distribution of the number of messages in a batch                                    |
	| consumer-group-join-total-<GroupID>       | counter    | Total count of consumer group join attempts                                          |
	| consumer-group-join-failed-<GroupID>      | counter    | Total count of consumer group join failures                                          |
	| consumer-group-sync-total-<GroupID>       | counter    | Total count of consumer group sync attempts                                          |
	| consumer-group-sync-failed-<GroupID>      | counter    | Total count of consumer group sync failures                                          |
	+-------------------------------------------+------------+--------------------------------------------------------------------------------------+
*/
package sarama

import (
	"io"
	"log"
)

var (
	// Logger is the instance of a StdLogger interface that Sarama writes connection
	// management events to. By default it is set to discard all log messages via ioutil.Discard,
	// but you can set it to redirect wherever you want.
	Logger StdLogger = log.New(io.Discard, "[Sarama] ", log.LstdFlags)

	// PanicHandler is called for recovering from panics spawned internally to the library (and thus
	// not recoverable by the caller's goroutine). Defaults to nil, which means panics are not recovered.
	PanicHandler func(interface{})

	// MaxRequestSize is the maximum size (in bytes) of any request that Sarama will attempt to send. Trying
	// to send a request larger than this will result in an PacketEncodingError. The default of 100 MiB is aligned
	// with Kafka's default `socket.request.max.bytes`, which is the largest request the broker will attempt
	// to process.
	MaxRequestSize int32 = 100 * 1024 * 1024

	// MaxResponseSize is the maximum size (in bytes) of any response that Sarama will attempt to parse. If
	// a broker returns a response message larger than this value, Sarama will return a PacketDecodingError to
	// protect the client from running out of memory. Please note that brokers do not have any natural limit on
	// the size of responses they send. In particular, they can send arbitrarily large fetch responses to consumers
	// (see https://issues.apache.org/jira/browse/KAFKA-2063).
	MaxResponseSize int32 = 100 * 1024 * 1024
)

// StdLogger is used to log error messages.
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type debugLogger struct{}

func (d *debugLogger) Print(v ...interface{}) {
	Logger.Print(v...)
}
func (d *debugLogger) Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}
func (d *debugLogger) Println(v ...interface{}) {
	Logger.Println(v...)
}

// DebugLogger is the instance of a StdLogger that Sarama writes more verbose
// debug information to. By default it is set to redirect all debug to the
// default Logger above, but you can optionally set it to another StdLogger
// instance to (e.g.,) discard debug information
var DebugLogger StdLogger = &debugLogger{}
