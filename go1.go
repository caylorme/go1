package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default:  %s", r.URL.Path)
}
func APIHandler(w http.ResponseWriter, r *http.Request) {
	Test := strings.Split(r.URL.Path, "/")
	switch strings.ToLower(Test[2]) {
	case "wow":
		fmt.Fprintf(w, "WOW!")
	default:
		Test2 := r.URL.Query()
		fmt.Printf("Test: %v\n", Test2.Get("test"))
		js, err := json.Marshal(Test2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func APIEncrypt(w http.ResponseWriter, r *http.Request) {
	k := []byte(r.URL.Query().Get("key"))
	text := r.URL.Query().Get("text")
	out := []byte(encrypt(k, text))
	w.Write(out)
}
func APIDecrypt(w http.ResponseWriter, r *http.Request) {
	k := []byte(r.URL.Query().Get("key"))
	text := r.URL.Query().Get("text")
	out := []byte(decrypt(k, text))
	w.Write(out)
}

func main() {
	fmt.Printf("Test main(): %v\n", decrypt([]byte("abcdefghijklmnsiwsznqjrheiwospew"), "mE8wRu4atz5IL3BubMahGpPRSHiN"))
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/API/", APIHandler)
	http.HandleFunc("/API/encrypt/", APIEncrypt)
	http.HandleFunc("/API/decrypt/", APIDecrypt)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
