package utils

import (
	"crypto/md5"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
)

var s *securecookie.SecureCookie

const saltPassowrd = "loooooooooooooooonggggggggggggHassssssssssheddddddddddddddSaaaaaaaaalt"

func init() {
	s = securecookie.New([]byte("very-secret"), []byte("aaaaaaaaaaaaaaaa"))

	// set age of cookie for 3 days
	s.MaxAge(86400 * 3)
}

//SetCookie
func SetCookie(cookieName, key string, value string) *http.Cookie {
	valueToEncrypt := map[string]string{
		key: value,
	}

	if encoded, err := s.Encode(cookieName, valueToEncrypt); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		return cookie
	}
	return nil
}

//GetCookie
func GetCookie(cookieName string, r *http.Request) map[string]string {
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if err = s.Decode(cookieName, cookie.Value, &value); err == nil {
			return value
		}
	}
	return nil
}

//encryptPassword: it helps in hasing password with mixed salt so
// it will be so hard to decrypt
func EncryptPassword(salt, password string) []byte {
	h := md5.New()
	io.WriteString(h, strings.Join([]string{salt, password}, ""))
	return h.Sum(nil)
}
