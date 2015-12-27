package main

import (
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
//	"mime"
	"mime/multipart"
	"net/textproto"
	"encoding/json"
	"strings"
	"os"
	"bytes"
    simplejson "github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
    ui "github.com/lshengjian/kms/webui"

)


var fileNameEscaper = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
func upload(c *echo.Context) error {
	mr, err := c.Request().MultipartReader()
	if err != nil {
		return err
	}

	// Read form field `name`
	part, err := mr.NextPart()
	if err != nil {
		return err
	}
	defer part.Close()
	b, err := ioutil.ReadAll(part)
	if err != nil {
		return err
	}
	name := string(b)


	// Read files
	i := 0
	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		defer part.Close()
        fname:=part.FileName()
        fmt.Println(fname)
        fullPathFilename:="uploads/"+fname
		file, err := os.Create(fullPathFilename)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, part); err != nil {
			return err
		}
		
		fh, _ := os.Open(fullPathFilename)
		url:="http://"+*web.filerServer+"/"+name+"/"+fname
		err=upload2FileServer(url, func(w io.Writer) (err error) {
		   _, err = io.Copy(w, fh)
		   return 
	        },fname)
		i++
		if(err!=nil) {
		  fmt.Println(err)
		}
	}


  args := make(map[string]interface{})  
  nodes,_ := listFileServer("http://"+*web.filerServer+"/"+name+"/")
  args["Author"] = name
  args["Master"]="http://"+*web.masterServer
  args["Files"] = nodes["Files"]
  return ui.FilesTpl.Execute(c.Response().Writer(), args)
	//return c.String(http.StatusOK, fmt.Sprintf("Thank You! %s , %d files uploaded successfully.",name,  i))
}


type UploadResult struct {
	Name  string `json:"name,omitempty"`
	Size  uint32 `json:"size,omitempty"`
	Error string `json:"error,omitempty"`
}
func listFileServer(url string) (map[string]interface{},error) {
    client:=&http.Client{Transport: &http.Transport{
		MaxIdleConnsPerHost: 1024,
	}}
	resp, err := client.Get(url)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()
	respBody, raErr := ioutil.ReadAll(resp.Body)
	if raErr != nil {
		return nil,raErr
	}
    js, _ := simplejson.NewJson(respBody)
    
	nodes, err2 := js.Map()
	//unmarshalErr := json.Unmarshal(resp_body, &ret)
   // fmt.Println(ret);
	return nodes, err2
}


func upload2FileServer(uploadUrl string, fillBufferFunction func(w io.Writer) error, filename string) (error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	h := make(textproto.MIMEHeader)
    fname:=fileNameEscaper.Replace(filename)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fname))
    //mtype := mime.TypeByExtension(strings.ToLower(filepath.Ext(filename)))
	
	
	file_writer, cp_err := body_writer.CreatePart(h)
	if cp_err != nil {
		return cp_err
	}
	if err := fillBufferFunction(file_writer); err != nil {
		return err
	}
	content_type := body_writer.FormDataContentType()
	if err := body_writer.Close(); err != nil {
		return err
	}
    client:=&http.Client{Transport: &http.Transport{
		MaxIdleConnsPerHost: 1024,
	}}
	resp, post_err := client.Post(uploadUrl, content_type, body_buf)
	if post_err != nil {
		return post_err
	}
	defer resp.Body.Close()
	resp_body, ra_err := ioutil.ReadAll(resp.Body)
	if ra_err != nil {
		return ra_err
	}
	var ret UploadResult
	unmarshal_err := json.Unmarshal(resp_body, &ret)
    fmt.Println(ret);
	return unmarshal_err
}

