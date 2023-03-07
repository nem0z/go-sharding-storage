package web

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nem0z/go-sharding-storage/types"
)

func HandleFile(ch_post chan *types.WrappedChan, ch_get chan *types.WrappedChan) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostFile(ch_post, w, r)
		case http.MethodGet:
			GetFile(ch_get, w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}
