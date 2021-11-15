package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/niroopreddym/HLSConversion/services"
)

//ShellToUse is the default cli type
const ShellToUse = "bash"

func main() {
	rdb := services.NewRedisInstance("127.0.0.1", "6389")

	// check if data is there in redis then read and give the data out else push teh data into redis
	if rdb.ContainsKey("test_video.mp4") {
		//read teh data from redis and convert the data to HLS
	} else {
		out, errout, err := Shellout(`../ffmpeg -i ../test_video.mp4 -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls ../OUTPUT/video/index.m3u8`)
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		fmt.Println("--- stdout ---")
		fmt.Println(out)
		fmt.Println("--- stderr ---")
		fmt.Println(errout)

		rdb.PlaceFFMPEGDataToRedis("../OUTPUT/video/", "test_video.mp4")
		res := cleanUpOutputFiles("../OUTPUT/video/")
		fmt.Println(res)
	}

}

//Shellout sets the default cli type to shell and stdout the response
func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func cleanUpOutputFiles(outPath string) bool {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := "rm -rf " + outPath + "/*.*"
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
