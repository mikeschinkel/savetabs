package restapi

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

func (a *API) catchPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var msg string

			reason := recover()
			if reason == nil {
				return
			}
			msg = fmt.Sprintf("%v", reason)
			rootDir := rootSourceDir()
			lines := strings.Split(string(debug.Stack()), "\n")
			for i := 0; i < len(lines); i++ {
				if !strings.HasPrefix(strings.TrimSpace(lines[i]), rootDir) {
					continue
				}
				msg = fmt.Sprintf("%s; %s", msg, strings.Join(lines[i-1:i+1], ""))
				break
			}
			slog.Error(msg)
			sendError(w, r, http.StatusInternalServerError, msg)
		}()
		next.ServeHTTP(w, r)
	})
}

func rootSourceDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filename))
}
