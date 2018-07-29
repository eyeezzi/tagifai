package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Image - an uploaded image
type Image struct {
	Name    string
	DataURL string
}

// Tag - a tag for an image
type Tag struct {
	Name  string
	Value float64
}

// uploadHandler - business logic to handle image submissions.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var img Image

	if err := json.NewDecoder(r.Body).Decode(&img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tags, err := getTags(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonRes, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(jsonRes))
}

// getTags - bare metal logic to access Clarifai API to get tags.
func getTags(img Image) (tags []Tag, err error) {

	apiKey, hasKey := os.LookupEnv("CLARIFAI_API_KEY")
	if !hasKey {
		log.Fatal("Missing environment variable: CLARIFAI_API_KEY")
	}

	url, hasURL := os.LookupEnv("IMG_ANALYSIS_ENDPOINT")
	if !hasURL {
		log.Fatal("Missing environment variable: IMG_ANALYSIS_ENDPOINT")
	}

	// construct request
	var fstr = `
	{
    "inputs": [
      {
        "data": {
          "image": {
						"base64": "%s"
          }
        }
      }
    ]
  }
	`
	dataStrs := strings.Split(img.DataURL, ",")
	b64str := dataStrs[len(dataStrs)-1]
	bodyStr := fmt.Sprintf(fstr, b64str)
	reqBody := []byte(bodyStr)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+apiKey)

	// make request with client
	tagClient := http.Client{}

	res, err := tagClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// handle response
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var resObj map[string]interface{}
	if err = json.Unmarshal(resBody, &resObj); err != nil {
		log.Fatal(err)
	}

	/*
		{
			"outputs": [
				{
					"data": {
						"concepts": [
							{}
						]
					}
				}
			]
		}
	*/
	concepts := resObj["outputs"].([]interface{})[0].(map[string]interface{})["data"].(map[string]interface{})["concepts"].([]interface{})

	for _, val := range concepts {
		var t Tag
		t.Name = val.(map[string]interface{})["name"].(string)
		t.Value = val.(map[string]interface{})["value"].(float64)

		tags = append(tags, t)
	}

	// fmt.Println(tags)
	return
}

// indexhandler - constructs and returns the only HTML page of the application
func indexHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"index.html"}
	ts := template.Must(template.ParseFiles(files...))
	ts.ExecuteTemplate(w, "layout", nil)
	return
}

func main() {
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tags", uploadHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on http://localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
