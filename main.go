package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/niroopreddym/HLSVideoStreaming/helpers"
	"github.com/niroopreddym/HLSVideoStreaming/services"
)

// ShellToUse is the default cli type
// const ShellToUse = "bash"
const ShellToUse = "sh"

func main() {
	rdbInstance := services.NewRedisInstance("127.0.0.1", "6379")
	fmt.Println(rdbInstance)

	//--------------create a temp directory in mount point for windwos machines----------
	tmp, err := os.MkdirTemp(os.TempDir(), "ffmpegtranscoding")
	if err != nil {
		log.Printf("error creating temp dir: %v\n", err)
	}

	defer os.RemoveAll(tmp)
	//--------------This above block will work for both windows and linux-----------------

	//--------------making use of io pipes on linux kernel----------
	// tmp, err := os.MkdirTemp(os.TempDir(), "ffmpegtranscoding")
	// if err != nil {
	// 	log.Printf("error creating temp dir: %v\n", err)
	// }

	// defer os.RemoveAll(tmp)
	//--------------This above block will work for both windows and linux-----------------

	//go routine to diwnload the S3 to fifo s3 > fifo

	//provide input as fifo to the below command and get the list in outLocation since in linux we dont have any other files int thsi location
	// its ok but in windows we are altering the code
	// outLocation := path.Join(dir, file.Name(), "index.m3u8")
	outputPath := filepath.Join(tmp, "index.m3u8")
	videoID := "test_video.mp4"
	cmd := exec.Command("ffmpeg",
		"-i", videoID,
		"-i", "TR.png", // Assuming your watermark image is a PNG file
		"-filter_complex", "overlay=W-w-10:H-h-10", // Adjust the position of the watermark
		"-profile:v", "baseline",
		"-level", "3.0",
		"-s", "640x360",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		`"`+outputPath+`"`,
	)

	cmdStr := cmd.String()
	fmt.Println(cmdStr)
	out, errout, err := Shellout(cmdStr)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)

	// since all the above operation happens in a go routine we need to monitor the go routine for err channel in a parllel for loop and report to the main.gof
	//use couple of go routines to push the data to redis
	ReadFiles(tmp+"/", "test_video.mp4", rdbInstance)
}

// ReadFiles reads all teh file sin the location
func ReadFiles(dir string, inputVideo string, rdb *services.Redis) {
	items, _ := os.ReadDir(dir)
	wg := sync.WaitGroup{}

	CleanMasterIndexerKeyForVideoEncodeFiles(rdb, inputVideo)
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		wg.Add(1)
		info, err := item.Info()
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}

		go MoveToRedis(dir, info, rdb, inputVideo, &wg)
	}

	wg.Wait()
}

// clean the master indexer key before pushing keys to list
func CleanMasterIndexerKeyForVideoEncodeFiles(rdb *services.Redis, inputVideo string) {
	if rdb.ContainsKey(inputVideo) {
		rdb.DeleteKey(inputVideo)
	}
}

// MoveToRedis moes data to redis
func MoveToRedis(dir string, item fs.FileInfo, rdb *services.Redis, inputVideo string, wg *sync.WaitGroup) {
	defer wg.Done()

	if item.IsDir() {
		subitems, _ := os.ReadDir(item.Name())
		for _, subitem := range subitems {
			if !subitem.IsDir() {
				fmt.Println(item.Name() + "/" + subitem.Name())
			}
		}
	} else {
		byteArrayItem := helpers.ConvertToByteArray(dir + item.Name())
		rdb.AddKeyValuePair(item.Name(), byteArrayItem)
		rdb.AddKeysToList(inputVideo, item.Name())
	}
}

// Shellout sets the default cli type to shell and stdout the response
func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
