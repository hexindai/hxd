// Copyright © 2017 Runrioter Wung <runrioter@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
	Long: `
	hxd tail --url=[连接]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(url) > 0 {
			resp, err := http.Head(url)
			if err != nil {
				log.Fatalln(err)
			}
			acceptRanges := resp.Header.Get("Accept-Ranges")
			contentLength := resp.Header.Get("Content-Length")
			if acceptRanges != "bytes" || contentLength == "" {
				log.Fatalln("该资源不支持基于http的tail, 请更改服务器配置。")
			}
			totalBytes, err := strconv.ParseInt(contentLength, 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%v\n", totalBytes)
			var startByte int64 = 0
			var endByte int64 = totalBytes
			if totalBytes-LastContentBytes >= 0 {
				startByte = totalBytes - LastContentBytes
			}

			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			req.Header.Add("Range", "bytes="+strconv.FormatInt(startByte, 10)+"-"+strconv.FormatInt(endByte, 10))
			respTail, err := client.Do(req)
			defer respTail.Body.Close()
			body, err := ioutil.ReadAll(respTail.Body)
			if err != nil {
				log.Fatalln(err)
			}
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
		} else {
			log.Fatalln("请输入url")
		}
	},
}

func init() {
	RootCmd.AddCommand(tailCmd)

	// Here you will define your flags and configuration settings.

	tailCmd.Flags().StringVar(&url, "url", "", "日志url连接")

}
