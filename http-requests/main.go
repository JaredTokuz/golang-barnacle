package main

import (
	"io/ioutil"

	"log"
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"

	"time"

	"io"
	"mime/multipart"
	"os"
)

func main() {
	// GetExample()
	// PostExample()
	// CustomRequest()
	PostFormJsonDecoderExample()
	//FileUpload()
}

func GetExample() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func PostExample() {
	requestBody, err := json.Marshal(map[string]string{
		"name": "as982k22jk",
		"email": "masnun@gmail.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func PostFormJsonDecoderExample() {
	formData := url.Values{
		"name": {"masnun"},
	}

	resp, err := http.PostForm("https://httpbin.org/post", formData)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result["form"])
}

func CustomRequest() {
	requestBody, err := json.Marshal(map[string]string{
		"name": "as982k22jk",
		"email": "masnun@gmail.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func FileUpload() {
	file, err := os.Open("name.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file_field", "name.txt")
	if err != nil {
		log.Fatalln(err)
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalln(err)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("normal_field")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		log.Fatalln(err)
	}

	// Completed adding the file and the fields, close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", "https://httpbin.org/post", &requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	// We need to set the content type from the writer it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(response.Body).Decode(&result)

	log.Println(result)
}