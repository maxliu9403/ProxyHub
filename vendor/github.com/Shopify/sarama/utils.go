package sarama

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
)

type none struct{}

// make []int32 sortable so we can sort partition numbers
type int32Slice []int32

func (slice int32Slice) Len() int {
	return len(slice)
}

func (slice int32Slice) Less(i, j int) bool {
	return slice[i] < slice[j]
}

func (slice int32Slice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func dupInt32Slice(input []int32) []int32 {
	ret := make([]int32, 0, len(input))
	ret = append(ret, input...)
	return ret
}

func withRecover(fn func()) {
	defer func() {
		handler := PanicHandler
		if handler != nil {
			if err := recover(); err != nil {
				handler(err)
			}
		}
	}()

	fn()
}

func safeAsyncClose(b *Broker) {
	tmp := b // local var prevents clobbering in goroutine
	go withRecover(func() {
		if connected, _ := tmp.Connected(); connected {
			if err := tmp.Close(); err != nil {
				Logger.Println("Error closing broker", tmp.ID(), ":", err)
			}
		}
	})
}

// Encoder is a simple interface for any type that can be encoded as an array of bytes
// in order to be sent as the key or value of a Kafka message. Length() is provided as an
// optimization, and must return the same as len() on the result of Encode().
type Encoder interface {
	Encode() ([]byte, error)
	Length() int
}

// make strings and byte slices encodable for convenience so they can be used as keys
// and/or values in kafka messages

// StringEncoder implements the Encoder interface for Go strings so that they can be used
// as the Key or Value in a ProducerMessage.
type StringEncoder string

func (s StringEncoder) Encode() ([]byte, error) {
	return []byte(s), nil
}

func (s StringEncoder) Length() int {
	return len(s)
}

// ByteEncoder implements the Encoder interface for Go byte slices so that they can be used
// as the Key or Value in a ProducerMessage.
type ByteEncoder []byte

func (b ByteEncoder) Encode() ([]byte, error) {
	return b, nil
}

func (b ByteEncoder) Length() int {
	return len(b)
}

// bufConn wraps a net.Conn with a buffer for reads to reduce the number of
// reads that trigger syscalls.
type bufConn struct {
	net.Conn
	buf *bufio.Reader
}

func newBufConn(conn net.Conn) *bufConn {
	return &bufConn{
		Conn: conn,
		buf:  bufio.NewReader(conn),
	}
}

func (bc *bufConn) Read(b []byte) (n int, err error) {
	return bc.buf.Read(b)
}

// KafkaVersion instances represent versions of the upstream Kafka broker.
type KafkaVersion struct {
	// it's a struct rather than just typing the array directly to make it opaque and stop people
	// generating their own arbitrary versions
	version [4]uint
}

func newKafkaVersion(major, minor, veryMinor, patch uint) KafkaVersion {
	return KafkaVersion{
		version: [4]uint{major, minor, veryMinor, patch},
	}
}

// IsAtLeast return true if and only if the version it is called on is
// greater than or equal to the version passed in:
//
//	V1.IsAtLeast(V2) // false
//	V2.IsAtLeast(V1) // true
func (v KafkaVersion) IsAtLeast(other KafkaVersion) bool {
	for i := range v.version {
		if v.version[i] > other.version[i] {
			return true
		} else if v.version[i] < other.version[i] {
			return false
		}
	}
	return true
}

// Effective constants defining the supported kafka versions.
var (
	V0_8_2_0  = newKafkaVersion(0, 8, 2, 0)
	V0_8_2_1  = newKafkaVersion(0, 8, 2, 1)
	V0_8_2_2  = newKafkaVersion(0, 8, 2, 2)
	V0_9_0_0  = newKafkaVersion(0, 9, 0, 0)
	V0_9_0_1  = newKafkaVersion(0, 9, 0, 1)
	V0_10_0_0 = newKafkaVersion(0, 10, 0, 0)
	V0_10_0_1 = newKafkaVersion(0, 10, 0, 1)
	V0_10_1_0 = newKafkaVersion(0, 10, 1, 0)
	V0_10_1_1 = newKafkaVersion(0, 10, 1, 1)
	V0_10_2_0 = newKafkaVersion(0, 10, 2, 0)
	V0_10_2_1 = newKafkaVersion(0, 10, 2, 1)
	V0_10_2_2 = newKafkaVersion(0, 10, 2, 2)
	V0_11_0_0 = newKafkaVersion(0, 11, 0, 0)
	V0_11_0_1 = newKafkaVersion(0, 11, 0, 1)
	V0_11_0_2 = newKafkaVersion(0, 11, 0, 2)
	V1_0_0_0  = newKafkaVersion(1, 0, 0, 0)
	V1_0_1_0  = newKafkaVersion(1, 0, 1, 0)
	V1_0_2_0  = newKafkaVersion(1, 0, 2, 0)
	V1_1_0_0  = newKafkaVersion(1, 1, 0, 0)
	V1_1_1_0  = newKafkaVersion(1, 1, 1, 0)
	V2_0_0_0  = newKafkaVersion(2, 0, 0, 0)
	V2_0_1_0  = newKafkaVersion(2, 0, 1, 0)
	V2_1_0_0  = newKafkaVersion(2, 1, 0, 0)
	V2_1_1_0  = newKafkaVersion(2, 1, 1, 0)
	V2_2_0_0  = newKafkaVersion(2, 2, 0, 0)
	V2_2_1_0  = newKafkaVersion(2, 2, 1, 0)
	V2_2_2_0  = newKafkaVersion(2, 2, 2, 0)
	V2_3_0_0  = newKafkaVersion(2, 3, 0, 0)
	V2_3_1_0  = newKafkaVersion(2, 3, 1, 0)
	V2_4_0_0  = newKafkaVersion(2, 4, 0, 0)
	V2_4_1_0  = newKafkaVersion(2, 4, 1, 0)
	V2_5_0_0  = newKafkaVersion(2, 5, 0, 0)
	V2_5_1_0  = newKafkaVersion(2, 5, 1, 0)
	V2_6_0_0  = newKafkaVersion(2, 6, 0, 0)
	V2_6_1_0  = newKafkaVersion(2, 6, 1, 0)
	V2_6_2_0  = newKafkaVersion(2, 6, 2, 0)
	V2_7_0_0  = newKafkaVersion(2, 7, 0, 0)
	V2_7_1_0  = newKafkaVersion(2, 7, 1, 0)
	V2_8_0_0  = newKafkaVersion(2, 8, 0, 0)
	V2_8_1_0  = newKafkaVersion(2, 8, 1, 0)
	V3_0_0_0  = newKafkaVersion(3, 0, 0, 0)

	SupportedVersions = []KafkaVersion{
		V0_8_2_0,
		V0_8_2_1,
		V0_8_2_2,
		V0_9_0_0,
		V0_9_0_1,
		V0_10_0_0,
		V0_10_0_1,
		V0_10_1_0,
		V0_10_1_1,
		V0_10_2_0,
		V0_10_2_1,
		V0_10_2_2,
		V0_11_0_0,
		V0_11_0_1,
		V0_11_0_2,
		V1_0_0_0,
		V1_0_1_0,
		V1_0_2_0,
		V1_1_0_0,
		V1_1_1_0,
		V2_0_0_0,
		V2_0_1_0,
		V2_1_0_0,
		V2_1_1_0,
		V2_2_0_0,
		V2_2_1_0,
		V2_2_2_0,
		V2_3_0_0,
		V2_3_1_0,
		V2_4_0_0,
		V2_4_1_0,
		V2_5_0_0,
		V2_5_1_0,
		V2_6_0_0,
		V2_6_1_0,
		V2_6_2_0,
		V2_7_0_0,
		V2_7_1_0,
		V2_8_0_0,
		V2_8_1_0,
		V3_0_0_0,
	}
	MinVersion     = V0_8_2_0
	MaxVersion     = V3_0_0_0
	DefaultVersion = V1_0_0_0
)

// ParseKafkaVersion parses and returns kafka version or error from a string
func ParseKafkaVersion(s string) (KafkaVersion, error) {
	if len(s) < 5 {
		return DefaultVersion, fmt.Errorf("invalid version `%s`", s)
	}
	var major, minor, veryMinor, patch uint
	var err error
	if s[0] == '0' {
		err = scanKafkaVersion(s, `^0\.\d+\.\d+\.\d+$`, "0.%d.%d.%d", [3]*uint{&minor, &veryMinor, &patch})
	} else {
		err = scanKafkaVersion(s, `^\d+\.\d+\.\d+$`, "%d.%d.%d", [3]*uint{&major, &minor, &veryMinor})
	}
	if err != nil {
		return DefaultVersion, err
	}
	return newKafkaVersion(major, minor, veryMinor, patch), nil
}

func scanKafkaVersion(s string, pattern string, format string, v [3]*uint) error {
	if !regexp.MustCompile(pattern).MatchString(s) {
		return fmt.Errorf("invalid version `%s`", s)
	}
	_, err := fmt.Sscanf(s, format, v[0], v[1], v[2])
	return err
}

func (v KafkaVersion) String() string {
	if v.version[0] == 0 {
		return fmt.Sprintf("0.%d.%d.%d", v.version[1], v.version[2], v.version[3])
	}

	return fmt.Sprintf("%d.%d.%d", v.version[0], v.version[1], v.version[2])
}
