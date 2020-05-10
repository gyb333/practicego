package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestJsonmarshal(t *testing.T)  {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalln(err)
	}
	var status = uint64(result["status"].(float64))
	fmt.Println("Status value: ", status)
}

func TestJsonDecode(t *testing.T) {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}
	var decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		log.Fatalln(err)
	}
	var status, _ = result["status"].(json.Number).Int64()
	fmt.Println("Status value: ", status)

}

func TestJsonStruct(t *testing.T)  {
	records := [][]byte{
		[]byte(`{"status":200, "tag":"one"}`),
		[]byte(`{"status":"ok", "tag":"two"}`),
	}
	for idx, record := range records {
		var result struct {
			StatusCode uint64
			StatusName string
			Status json.RawMessage `json:"status"`
			Tag string `json:"tag"`
		}
		err := json.NewDecoder(bytes.NewReader(record)).Decode(&result)
		if err!=nil{
			panic(err)
		}

		err = json.Unmarshal(result.Status, &result.StatusName)
		if err != nil {
			err =json.Unmarshal(result.Status,&result.StatusCode)
		}
		fmt.Println(idx,result )
	}
}
