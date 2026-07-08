package cmdutil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// ReadBody resolves the body text from --body and --body-file flags.
// body and bodyFile are the raw flag values; stdin is used when bodyFile == "-".
// It returns an error when both body and bodyFile are set.
func ReadBody(body, bodyFile string, stdin io.Reader) (string, error) {
	if body != "" && bodyFile != "" {
		return "", fmt.Errorf("cannot use both --body and --body-file")
	}

	if body != "" {
		if err := ScanContentForSecrets(body); err != nil {
			return "", err
		}
		return body, nil
	}

	if bodyFile != "" {
		if bodyFile == "-" {
			bodyText, err := ReadTextFromFlag(stdin, "--body-file")
			if err != nil {
				return "", fmt.Errorf("failed to read from stdin: %w", err)
			}
			return strings.TrimSpace(bodyText), nil
		}

		// Read from file
		content, err := ReadTextFile(bodyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", bodyFile, err)
		}
		return strings.TrimSpace(content), nil
	}

	return "", nil
}

var utf8BOM = []byte{0xef, 0xbb, 0xbf}
var utf16LEBOM = []byte{0xff, 0xfe}
var utf16BEBOM = []byte{0xfe, 0xff}

// ErrLossyPowerShellStdin is returned when Windows PowerShell appears to have
// replaced non-ASCII stdin text with question marks before the CLI could read it.
var ErrLossyPowerShellStdin = errors.New("lossy Windows PowerShell stdin")

// ReadTextFile reads a user-provided text file and strips a UTF-8 BOM when present.
func ReadTextFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	text := DecodeUserText(content)
	if err := ScanContentForSecrets(text); err != nil {
		return "", err
	}
	return text, nil
}

// ReadText reads user-provided text from a stream.
func ReadText(r io.Reader) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return DecodeUserText(content), nil
}

// ReadTextFromFlag reads stdin text for an explicit file flag such as
// --body-file - or --comment-file - and rejects input that Windows PowerShell
// appears to have already corrupted before the CLI could decode it.
func ReadTextFromFlag(r io.Reader, flagName string) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	text := DecodeUserText(content)
	if isLikelyLossyPowerShellStdin(content, text, runtime.GOOS, os.Getenv("GITCODE_CLI_ALLOW_LOSSY_STDIN") != "") {
		return "", newLossyPowerShellStdinError(flagName)
	}
	if err := ScanContentForSecrets(text); err != nil {
		return "", err
	}
	return text, nil
}

// DecodeUserText decodes text accepted from user-facing files or stdin.
func DecodeUserText(content []byte) string {
	content = bytes.TrimPrefix(content, utf8BOM)
	if bytes.HasPrefix(content, utf16LEBOM) {
		return decodeUTF16(content[len(utf16LEBOM):], binary.LittleEndian)
	}
	if bytes.HasPrefix(content, utf16BEBOM) {
		return decodeUTF16(content[len(utf16BEBOM):], binary.BigEndian)
	}
	if utf8.Valid(content) {
		return string(content)
	}
	if decoded, err := simplifiedchinese.GB18030.NewDecoder().Bytes(content); err == nil && utf8.Valid(decoded) {
		return string(decoded)
	}
	return string(bytes.ToValidUTF8(content, []byte("\uFFFD")))
}

func decodeUTF16(content []byte, order binary.ByteOrder) string {
	if len(content)%2 == 1 {
		content = content[:len(content)-1]
	}
	u16 := make([]uint16, 0, len(content)/2)
	for len(content) >= 2 {
		u16 = append(u16, order.Uint16(content[:2]))
		content = content[2:]
	}
	return string(utf16.Decode(u16))
}

func isLikelyLossyPowerShellStdin(raw []byte, text, goos string, allowLossy bool) bool {
	if allowLossy || goos != "windows" || len(raw) == 0 {
		return false
	}
	if !isASCII(raw) {
		return false
	}
	return hasQuestionMarkRun(text, 3)
}

func isASCII(raw []byte) bool {
	for _, b := range raw {
		if b >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

func hasQuestionMarkRun(text string, minRun int) bool {
	return strings.Contains(text, strings.Repeat("?", minRun))
}

func newLossyPowerShellStdinError(flagName string) error {
	if flagName == "" {
		flagName = "--body-file"
	}
	command := "gitcode issue create -R owner/repo --title \"标题\""
	fileName := "body.md"
	if flagName == "--comment-file" {
		command = "gitcode pr review 1 -R owner/repo"
		fileName = "comment.md"
	}
	return fmt.Errorf("%w: Windows PowerShell seems to have replaced Chinese/non-ASCII stdin text with ??? before GitCode CLI received it.\n\nCorrect usage / 正确用法:\n  $OutputEncoding = [System.Text.UTF8Encoding]::new($false)\n  \"中文正文\" | %s %s -\n\nMore stable / 更稳妥:\n  Set-Content -Path %s -Value \"中文正文\" -Encoding UTF8\n  %s %s %s\n\nUse the same %s flag with the command you are running. If the literal question marks are intentional, set GITCODE_CLI_ALLOW_LOSSY_STDIN=1.", ErrLossyPowerShellStdin, command, flagName, fileName, command, flagName, fileName, flagName)
}
