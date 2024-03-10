package middlewares

import (
	"fmt"
	"github.com/Jeff-Rdg/webhook_discord/messenger"
	"github.com/spf13/viper"
	"net/http"
)

type ResponseCapturer struct {
	http.ResponseWriter
	statusCode int
	body       string
	requestURI string
	requestURL string
}

func (rc *ResponseCapturer) WriteHeader(statusCode int) {
	rc.statusCode = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}

func (rc *ResponseCapturer) Write(b []byte) (int, error) {
	rc.body = string(b)
	return rc.ResponseWriter.Write(b)
}

func getRequestURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
}

func LogException(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &ResponseCapturer{ResponseWriter: w,
			requestURI: r.RequestURI,
			requestURL: getRequestURL(r),
		}
		next.ServeHTTP(rc, r)

		status := rc.statusCode

		if status > 300 {
			messenger.SendLog(messenger.SenderLog{
				Enviroment:      viper.GetString("ENVIROMENT"),
				Response:        rc.body,
				Uri:             rc.requestURI,
				UriUrl:          rc.requestURL,
				ApplicationName: viper.GetString("APP_NAME"),
				WebhookUrl:      viper.GetString("WEBHOOK_URL"),
			})
		}

	})

}
