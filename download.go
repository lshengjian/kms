package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"github.com/lshengjian/kms/operation"
	"github.com/lshengjian/kms/util"
)

var (
	d DownloadOptions
)

type DownloadOptions struct {
	server *string
	dir    *string
}

func init() {
	cmdDownload.ExecuteFunc = runDownload // break init cycle
	d.server = cmdDownload.Flag.String("server", "localhost:9333", "SeaweedFS master location")
	d.dir = cmdDownload.Flag.String("dir", "downloads", "Download the whole folder recursively if specified.")
}

var cmdDownload = &Command{
	UsageLine: "download -server=localhost:9333 -dir=downloads fid1 [fid2 fid3 ...]",
	Short:     "download files by file id",
	Long: `download files by file id.

  Usually you just need to use curl to lookup the file's volume server, and then download them directly.
  This download tool combine the two steps into one.

  `,
}

func runDownload(cmd *Command, args []string) bool {
	for _, fid := range args {
		filename, content, e := fetchFileId(*d.server, fid)
		if e != nil {
			fmt.Println("Fetch Error:", e)
			continue
		}
		if filename == "" {
			filename = fid
		}
		if strings.HasSuffix(filename, "-list") {
			filename = filename[0 : len(filename)-len("-list")]
			fids := strings.Split(string(content), "\n")
			f, err := os.OpenFile(path.Join(*d.dir, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				fmt.Println("File Creation Error:", e)
				continue
			}
			defer f.Close()
			for _, partId := range fids {
				var n int
				_, part, err := fetchFileId(*d.server, partId)
				if err == nil {
					n, err = f.Write(part)
				}
				if err == nil && n < len(part) {
					err = io.ErrShortWrite
				}
				if err != nil {
					fmt.Println("File Write Error:", err)
					break
				}
			}
		} else {
			ioutil.WriteFile(path.Join(*d.dir, filename), content, os.ModePerm)
		}
	}
	return true
}

func fetchFileId(server string, fileId string) (filename string, content []byte, e error) {
	fileUrl, lookupError := operation.LookupFileId(server, fileId)
	if lookupError != nil {
		return "", nil, lookupError
	}
	filename, content, e = util.DownloadUrl(fileUrl)
	return
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	f.Close()
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}
