package middleware

import (
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            //deferで仕込む
            defer func() {
                if err := recover(); err != nil {

                    log.Printf("[PANIC RECOVERED] err=%v\n", err)

                    http.Error(w,"Internal Server Error", http.StatusInternalServerError)
                }
            }()
        h.ServeHTTP(w, r)

    }

    return http.HandlerFunc(fn)
}