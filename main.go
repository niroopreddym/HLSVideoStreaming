package main

import (
	"fmt"

	"github.com/niroopreddym/HLSConversion/services"
)

func main() {
	rdbInstance := services.NewRedisInstance("127.0.0.1", "6389")
	fmt.Println(rdbInstance)

	rdbInstance.PlaceFFMPEGDataToRedis("OUTPUT/video/", "test_video.mp4")
	//spin up a mux api to fetch list of vieo chunks
}
