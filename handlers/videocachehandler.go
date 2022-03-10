package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/niroopreddym/HLSConversion/services"
)

//VideoCache is the Cache Handler
type VideoCache struct {
	RedisDB   services.Redis
	muxRouter *mux.Router
}

//NewVideoCacheHandler returns a new instance of the VideoCacheHandler
func NewVideoCacheHandler(mux *mux.Router) *VideoCache {
	return &VideoCache{
		RedisDB:   *services.NewRedisInstance("127.0.0.1", "6389"),
		muxRouter: mux,
	}
}

//GetPlaylistInfo retrives the Meta File
func (handler *VideoCache) GetPlaylistInfo(w http.ResponseWriter, r *http.Request) {
	condition := false
	if condition {
		params := mux.Vars(r)
		videoID := params["id"]
		cmd := handler.RedisDB.RedisClient.LRange(string(videoID), 0, 20)
		val := cmd.Val()
		indexFile := val[0]
		value := handler.RedisDB.GetValueByKey(indexFile)
		//alter the m3u8File and append the prefix to the ts files
		alteredString := alterTheM3U8Data(value, videoID)
		b := bytes.NewBuffer([]byte(alteredString))

		// stream straight to client(browser)
		w.Header().Set("Content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if _, err := b.WriteTo(w); err != nil {
			fmt.Fprintf(w, "%s", err)
		}
	} else {
		handler.muxRouter.PathPrefix("/asset/video/").Handler(http.StripPrefix("/asset/video/", http.FileServer(http.Dir("../OUTPUT/video"))))
	}

	return
}

func alterTheM3U8Data(value string, videoID string) string {
	val := strings.Replace(value, "index", "/asset/video/"+videoID+"/index", -1)
	return val
}

//StreamFileSegments streams the segment bytes to the Client
func (handler *VideoCache) StreamFileSegments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// videoID := params["id"]
	segmentID := params["segment_id"]
	value := handler.RedisDB.GetValueByKey(segmentID)
	b := bytes.NewBuffer([]byte(value))

	// stream straight to client(browser)
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if _, err := b.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	return
}

// func responseController(w http.ResponseWriter, code int, payload interface{}) {
// 	response, _ := json.Marshal(payload)
// 	// w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Content-Type", "application/x-mpegURL")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }

func responseController(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}
