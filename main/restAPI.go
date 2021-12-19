package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"time"
)

// Struct for response recieved from api.random.org
type ResponseRandom struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Random struct {
			Data           []int  `json:"data"`
			CompletionTime string `json:"completionTime"`
		} `json:"random"`
		BitsUsed      int `json:"bitsUsed"`
		BitsLeft      int `json:"bitsLeft"`
		RequestsLeft  int `json:"requestsLeft"`
		AdvisoryDelay int `json:"advisoryDelay"`
	} `json:"result"`
	ID int `json:"id"`
}
// Struct for response of our RESTapi
type Response struct {
	Stddev int   `json:"stddev"`
	Data   []int `json:"data"`
}
// Struct for request which we are gonna send to the api.random.org
type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		APIKey                    string      `json:"apiKey"`
		N                         int         `json:"n"`
		Min                       int         `json:"min"`
		Max                       int         `json:"max"`
		Replacement               bool        `json:"replacement"`
		Base                      int         `json:"base"`
		PregeneratedRandomization interface{} `json:"pregeneratedRandomization"`
	} `json:"params"`
	ID int `json:"id"`
}

// Struct to store all responses together
type combinedRespones []Response

var responses = combinedRespones{}

func getResponse(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	select {
    case <-time.After(2 * time.Second):
		url := "https://api.random.org/json-rpc/4/invoke"

		vars := r.URL.Query()

		_len, err := strconv.Atoi(vars["length"][0])
		_req, err := strconv.Atoi(vars["request"][0])

		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		// Request object for api.random.org
		req := Request{
			Jsonrpc: "2.0",
			Method:  "generateIntegers",
		}
		// Params for our requests
		req.Params.APIKey = "3038998a-88c2-445a-aed3-7398e03d32f2"
		req.Params.Min = 1
		req.Params.N = _len
		req.Params.Max = 10
		req.Params.Replacement = true
		req.Params.Base = 10


		var c_response[]int
		for r := 0; r < _req; r++ {

			client := &http.Client{Timeout: 15 * time.Second}

			jsonReq, _ := json.Marshal(req)
			request, error := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
			if error != nil {
				panic(error)
			}
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			response, error := client.Do(request)
			if error != nil {
				panic(error)
			} else if response.StatusCode != 200 {
				panic(response.StatusCode)
			}

			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			var responseObj ResponseRandom
			json.Unmarshal(body, &responseObj)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			if responseObj.Result.Random.Data == nil {
				panic("Invalid length argument! Value must be an integer!")
			} else {
				res := Response{
					Stddev: 1,
					Data:   responseObj.Result.Random.Data,
				}
				responses = append(responses, res)
				c_response = append(c_response, responses[r].Data...)
				
			}
		}
		// Object of item in our json resposne, which contains combined data
		summed_res := Response{
			Stddev: 1,
			Data: c_response,
		}
		responses = append(responses, summed_res)
		json.NewEncoder(w).Encode(responses)
		responses = nil

	// Context error handler
	case <-ctx.Done():
		err := ctx.Err()
        fmt.Println("server:", err)
        internalError := http.StatusInternalServerError
        http.Error(w, err.Error(), internalError)
	case <-time.After(15*time.Second):
		err := ctx.Err()
		fmt.Println("server:", err)
		timeoutError := http.ErrHandlerTimeout
		fmt.Println("timeout:", timeoutError)
    }
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/random/mean", getResponse)
	log.Fatal(http.ListenAndServe(":8000", router))
}
