package loader

import (
	"bytes"
	"fmt"
	"golang.org/x/text/transform"
	"math/rand"
	"regexp"
	"time"

	"golang.org/x/text/encoding/unicode"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

func ObfuscateStrings(b []byte, blacklist []string) []byte{
	fmt.Println("Replapcing %d keywords", len(blacklist))
	for _, word := range blacklist {
		b = ReplaceWord(b, word)
	}
	return b
}

func ReplaceWord(b []byte, word string) []byte{
	newWord := shuffle(word)
	re := regexp.MustCompile("(?i)" + utf16LeStr(word))
	b = re.ReplaceAll(b, utf16Le(newWord))
	re2 := regexp.MustCompile("(?i)" + word)
	b = re2.ReplaceAll(b, []byte(newWord))
	fmt.Println("Replacing %s with %s", word, newWord)
	return b
}


func utf16Le(s string) []byte {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	var buf bytes.Buffer
	t := transform.NewWriter(&buf, enc)
	t.Write([]byte(s))
	return buf.Bytes()
}

func utf16LeStr(s string) string {
	return string(utf16Le(s))
}


func shuffle(in string) string {
	inRune := []rune(in)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
