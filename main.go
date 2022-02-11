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

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)
	var seedstr string
	var formatBase64, formatBinary, omitNewline bool

	cmd := &cobra.Command{
		Use:   "rand LENGTH_BYTES",
		Short: "Print random bytes from a secure source to stdout.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			n, err := strconv.ParseUint(args[0], 10, 32)
			fatal(err, "argument must be a positive 32bit integer")
			bs := make([]byte, int(n))
			var read func([]byte) (int, error) = rand.Read
			sourcelabel := "secure"

			if seedstr != "" {
				seed, err := strconv.ParseInt(seedstr, 10, 64)
				fatal(err, "invalid seed value")
				mathrand.Seed(seed)
				sourcelabel = "insecure"
				read = mathrand.Read
			}
			_, err = read(bs)
			fatal(err, "failed to read random bytes from %v source", sourcelabel)

			trailingNewline := "\n"
			if omitNewline {
				trailingNewline = ""
			}

			switch true {
			case formatBase64 && formatBinary:
				fatal(fmt.Errorf(`"formatBase64", "formatBinary"`), "incompatible flags")
			case formatBase64:
				fmt.Print(base64.StdEncoding.EncodeToString(bs) + trailingNewline)
			case formatBinary:
				_, err = os.Stdout.Write(bs)
				fatal(err, "failed to write unformatted bytes to stdout")
			default:
				fmt.Print(hex.EncodeToString(bs) + trailingNewline)
			}
		},

		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	flags := cmd.Flags()
	flags.StringVarP(&seedstr, "seed", "s", "", "seed value as a decimal 64bit integer using an insecure random source")
	flags.BoolVarP(&formatBase64, "base64", "a", false, "print random bytes encoded as Base64")
	flags.BoolVarP(&formatBinary, "binary", "b", false, "print random bytes directly without formatting")
	flags.BoolVarP(&omitNewline, "omit-newline", "n", false, "do not print the trailing newline character")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func fatal(err error, message string, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			message = fmt.Sprintf(message, args...)
		}
		log.Fatalf("FATAL: %v: %v", message, err)
	}
}
