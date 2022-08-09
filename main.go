package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var wg sync.WaitGroup

type PromotionRecord struct {
	id              string
	price           float64
	expiration_date string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "YOU ARE WELCOME!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {

	/*f, err := os.Open("promotions.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	promotionList := CreatePromotionList(data)

	// print the array
	fmt.Printf("%+v\n", promotionList)
	json_data, err := json.Marshal(promotionList)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "%v", string(json_data))
	//create json file
	json_file, err := os.Create("promotionList.json")
	if err != nil {
		fmt.Println(err)
	}
	defer json_file.Close()

	json_file.Write(json_data)
	json_file.Close()
	//print json data
	*/
}

func HTTPhandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ch := make(chan []string)
	wg.Add(1)
	go csvReader(id, w, ch) //r.URL.Path[1:]
	resultRow := <-ch       //go csvReader(r.URL.Path[1:], w, ch)
	sizeResult := len(resultRow)
	wg.Wait()
	//fmt.Fprintf(w, "%v  size: %v", resultRow, len(ch))

	if sizeResult > 0 {
		fmt.Fprintf(w, "ID num: %v Price: %v Time: %v", resultRow[0], resultRow[1], resultRow[2]) // fmt.Fprintf(w, "%v ", resultRow[0])
	} else {
		fmt.Fprintf(w, "No Matching row. ")
	}
}

func main() {
	//var web webServer
	http.HandleFunc("/", HTTPhandler)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	//myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/promotion/{id}", HTTPhandler) //returnSingleArticle
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func CSVFileCreator() {
	data := [][]string{
		//{"id", "price", "expiration_date"},
		{"d018ef0b-dbd9-48f1-ac1a-eb4d90e57118", "100", "2018-08-04 05:32:31 +0200 CEST"},
	}

	// create a file
	file, err := os.Create("promotions_sample.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)

	defer writer.Flush()

	// write all rows at once
	writer.WriteAll(data)
}

func csvReader(data string, w http.ResponseWriter, ch chan<- []string) {
	f, err := os.Open("promotions.csv")
	var recordResult []string
	if err != nil {
		log.Fatal(err)
	}

	// initialize csv reader
	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if record[0] == data {
			//fmt.Printf("ID: %s ", record)
			recordResult = append(recordResult, record...)
			fmt.Printf("ID: %s size: %v ", recordResult, len(recordResult))
			json_data, err := json.Marshal(recordResult)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//print json data
			fmt.Println(string(json_data))
			break
		}
	}
	//Workgroup Done
	wg.Done()
	ch <- recordResult
}

func CSVReaderAll() {

}

func CreatePromotionList(data [][]string) []PromotionRecord {
	var promotionList []PromotionRecord
	for i, line := range data {
		if i > 0 { // omit header line
			var rec PromotionRecord
			for j, field := range line {
				if j == 0 {
					rec.id = field
				} else if j == 1 {
					rec.price, _ = strconv.ParseFloat(field, 64)
				} else if j == 2 {
					rec.expiration_date = field
				}
			}
			promotionList = append(promotionList, rec)
		}
	}
	return promotionList
}
