package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	mathrand "math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var formatBase64, formatBinary, formatPassword, omitNewline, formatUUID bool
var omit string

func main() {
	log.SetFlags(0)
	var seedstr string

	cmd := &cobra.Command{
		Use:   "rand LENGTH_BYTES",
		Short: "Print random bytes from a secure source to stdout.",
		Run: func(cmd *cobra.Command, args []string) {
			var randSize int
			if formatUUID {
				randSize = 16
			} else if len(args) != 1 {
				n, err := strconv.ParseUint(args[0], 10, 32)
				fatal(err, "argument must be a positive 32bit integer")
				randSize = int(n)
			} else {
				cmd.Help()
				os.Exit(1)
			}
			bs := make([]byte, randSize)
			read := rand.Read
			sourcelabel := "secure"

			if seedstr != "" {
				seed, err := strconv.ParseInt(seedstr, 10, 64)
				fatal(err, "invalid seed value")
				mathrand.Seed(seed)
				sourcelabel = "insecure"
				read = mathrand.Read
			}
			_, err := read(bs)
			fatal(err, "failed to read random bytes from %v source", sourcelabel)

			trailingNewline := "\n"
			if omitNewline {
				trailingNewline = ""
			}

			formats := listFormats()

			switch true {
			case len(formats) > 1:
				fatal(fmt.Errorf(strings.Join(formats, ", ")), "incompatible flags")
			case formatBase64:
				fmt.Print(base64.StdEncoding.EncodeToString(bs) + trailingNewline)
			case formatBinary:
				_, err = os.Stdout.Write(bs)
				fatal(err, "failed to write unformatted bytes to stdout")
			case formatPassword:
				fmt.Print(encodePassword(bs) + trailingNewline)
			case formatUUID:
				bs[6] = (bs[6] & 0x0f) | 0x40
				bs[8] = (bs[8] & 0x3f) | 0x80
				fmt.Printf("%x-%x-%x-%x-%x%v", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:], trailingNewline)
			default:
				fmt.Print(hex.EncodeToString(bs) + trailingNewline)
			}
		},

		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	flags := cmd.Flags()
	flags.StringVarP(&seedstr, "seed", "s", "", "use an insecure random source with seed integer")
	flags.BoolVarP(&formatBase64, "base64", "a", false, "print random bytes encoded as base64")
	flags.BoolVarP(&formatBinary, "binary", "b", false, "print random bytes directly without formatting")
	flags.BoolVarP(&formatPassword, "password", "p", false, "print a suitable password")
	flags.StringVarP(&omit, "omit", "o", "", "omit the listed characters from generated passwords")
	flags.BoolVarP(&omitNewline, "omit-newline", "n", false, "do not print the trailing newline character")
	flags.BoolVarP(&formatUUID, "uuid", "u", false, "print a random (v4) UUID")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func listFormats() []string {
	fm := map[string]bool{
		"base64":   formatBase64,
		"binary":   formatBinary,
		"password": formatPassword,
		"uuid":     formatUUID,
	}
	fs := make([]string, 0, len(fm))
	for k, v := range fm {
		if v {
			fs = append(fs, fmt.Sprintf(`"--%v"`, k))
		}
	}
	return fs
}

const passwordChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-=!@#$%^&*()_+[]\\{}|;:,./<>?~"

func encodePassword(bs []byte) string {
	pass := make([]byte, len(bs))
	passwordBytes := []byte(passwordChars)
	for i, b := range bs {
		c := passwordBytes[int(b)%len(passwordBytes)]
		for strings.ContainsAny(string([]byte{c}), omit) {
			c = passwordBytes[mathrand.Intn(len(passwordBytes))]
		}
		pass[i] = c
	}
	return string(pass)
}

func fatal(err error, message string, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			message = fmt.Sprintf(message, args...)
		}
		log.Fatalf("FATAL: %v: %v", message, err)
	}
}
