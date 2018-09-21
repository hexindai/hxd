package cmd

import (
	"log"
	"strconv"
	"time"

	"net/http"

	"fmt"

	"io/ioutil"

	"github.com/spf13/cobra"
)

var url string

// LastContentBytes 预留的信息
const LastContentBytes = 5000

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "动态查看日志文件",
	Run: func(cmd *cobra.Command, args []string) {

		if len(url) == 0 {
			log.Fatalln("url required")
		}

		resp, err := http.Head(url)
		if err != nil {
			log.Fatalln(err)
		}

		acceptRanges := resp.Header.Get("Accept-Ranges")
		contentLength := resp.Header.Get("Content-Length")

		if resp.StatusCode == http.StatusNotFound {
			log.Fatalln(errNotFound)
		}

		if acceptRanges != "bytes" || contentLength == "" {
			log.Fatalln(errNotSupported)
		}

		totalBytes, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}

		startByte := int64(0)
		endByte := totalBytes
		if totalBytes-LastContentBytes >= 0 {
			startByte = totalBytes - LastContentBytes
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Range", "bytes="+strconv.FormatInt(startByte, 10)+"-"+strconv.FormatInt(endByte, 10))
		respTail, err := client.Do(req)

		body, err := ioutil.ReadAll(respTail.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer respTail.Body.Close()

		startByte = endByte
		fmt.Printf("%s", body)

		for {
			time.Sleep(2 * time.Second)
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)

			req.Header.Add("Range", "bytes="+strconv.FormatInt(startByte, 10)+"-")
			respTail, err := client.Do(req)
			if respTail.StatusCode == http.StatusPartialContent {
				offset := respTail.Header.Get("Content-Length")
				offsetBytes, err := strconv.ParseInt(offset, 10, 64)
				if err != nil {
					log.Fatalln(err)
				}
				startByte += offsetBytes
			} else if respTail.StatusCode == http.StatusRequestedRangeNotSatisfiable {
				resp, err := http.Head(url)
				if err != nil {
					log.Fatalln(err)
				}
				cL := resp.Header.Get("Content-Length")
				newByte, err := strconv.ParseInt(cL, 10, 64)
				if newByte < startByte {
					startByte = 0
				}
				continue
			}
			defer respTail.Body.Close()
			body, err := ioutil.ReadAll(respTail.Body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%s", body)
		}
	},
}

func init() {
	hxd.AddCommand(tailCmd)
	tailCmd.Flags().StringVar(&url, "url", "", "http url of log file")
	tailCmd.MarkFlagRequired("url")
}
