package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/niroopreddym/HLSConversion/helpers"
	"github.com/niroopreddym/HLSConversion/services"
)

//ShellToUse is the default cli type
const ShellToUse = "bash"

func main() {
	rdbInstance := services.NewRedisInstance("127.0.0.1", "6379")
	fmt.Println(rdbInstance)

	// rdbInstance.PlaceFFMPEGDataToRedis("OUTPUT/video/", "test_video.mp4")
	// //spin up a mux api to fetch list of vieo chunks

	//create a temp directory in mount point

	dir, err := ioutil.TempDir("./mountpoint/", "hls.XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)
	//create a fifo in dir in windows let it be tempfile
	file, err := ioutil.TempFile(dir, "prefix")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file.Name())
	//go routine to diwnload the S3 to fifo s3 > fifo

	//provide input as fifo to the below command and get the list in outLocation since in linux we dont have any other files int thsi location
	// its ok but in windows we are altering the code
	// outLocation := path.Join(dir, file.Name(), "index.m3u8")
	out, errout, err = Shellout(`ffmpeg -i test_video.mp4 -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls OUTPUT/video/index.m3u8`)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)

	// since all the above operation happens in a go routine we need to monitor the go routine for err channel in a parllel for loop and report to the main.go
	//use couple of go routines to push the data to redis
	ReadFiles("OUTPUT/video/", "test_video.mp4", rdbInstance)

}

//ReadFiles reads all teh file sin the location
func ReadFiles(dir string, inputVideo string, rdb *services.Redis) {
	items, _ := ioutil.ReadDir(dir)
	wg := sync.WaitGroup{}

	for _, item := range items {
		wg.Add(1)
		go MoveToRedis(dir, item, rdb, inputVideo, &wg)
	}

	wg.Wait()
}

//MoveToRedis moes data to redis
func MoveToRedis(dir string, item fs.FileInfo, rdb *services.Redis, inputVideo string, wg *sync.WaitGroup) {
	defer wg.Done()
	if item.IsDir() {
		subitems, _ := ioutil.ReadDir(item.Name())
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
