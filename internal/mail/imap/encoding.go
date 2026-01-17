package imap

import (
	"bytes"
	"encoding/hex"
	"io"
	"strings"

	"github.com/emersion/go-message/charset"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

func init() {
	// Register common charsets that aren't registered by default
	charset.RegisterEncoding("ascii", unicode.UTF8)
	charset.RegisterEncoding("us-ascii", unicode.UTF8)
	charset.RegisterEncoding("iso-8859-1", charmap.ISO8859_1)
	charset.RegisterEncoding("iso-8859-2", charmap.ISO8859_2)
	charset.RegisterEncoding("iso-8859-15", charmap.ISO8859_15)
	charset.RegisterEncoding("windows-1252", charmap.Windows1252)
}

// decodeTransferEncoding decodes based on Content-Transfer-Encoding header
func decodeTransferEncoding(body []byte, encoding string) string {
	encoding = strings.ToLower(strings.TrimSpace(encoding))

	switch encoding {
	case "quoted-printable":
		return decodeQuotedPrintable(body)

	case "base64":
		// Base64 is usually auto-decoded by mail.CreateReader
		// But decode manually if needed
		return string(body)

	case "7bit", "8bit", "binary", "":
		return string(body)

	default:
		return string(body)
	}
}

// decodeQuotedPrintable manually decodes quoted-printable (no buffer limits)
func decodeQuotedPrintable(input []byte) string {
	var result bytes.Buffer
	result.Grow(len(input)) // Pre-allocate

	i := 0
	for i < len(input) {
		c := input[i]

		if c == '=' {
			// Check for soft line break (=\r\n or =\n)
			if i+1 < len(input) && input[i+1] == '\n' {
				// Soft line break: =\n
				i += 2
				continue
			}
			if i+2 < len(input) && input[i+1] == '\r' && input[i+2] == '\n' {
				// Soft line break: =\r\n
				i += 3
				continue
			}

			// Try to decode hex (e.g., =3D → =)
			if i+2 < len(input) {
				hexStr := string(input[i+1 : i+3])
				decoded, err := hex.DecodeString(hexStr)
				if err == nil && len(decoded) == 1 {
					result.WriteByte(decoded[0])
					i += 3
					continue
				}
			}

			// Can't decode - write '=' as is
			result.WriteByte(c)
			i++
		} else {
			// Normal character
			result.WriteByte(c)
			i++
		}
	}

	return result.String()
}

// convertCharset converts text from specified charset to UTF-8
func convertCharset(text string, fromCharset string) string {
	fromCharset = strings.ToLower(fromCharset)

	// Already UTF-8 or ASCII
	if fromCharset == "utf-8" || fromCharset == "us-ascii" || fromCharset == "ascii" {
		return text
	}

	// Try to convert using charset package
	reader, err := charset.Reader(fromCharset, strings.NewReader(text))
	if err != nil {
		// Can't convert - return as is
		return text
	}

	converted, err := io.ReadAll(reader)
	if err != nil {
		// Conversion failed - return original
		return text
	}

	return string(converted)
}
