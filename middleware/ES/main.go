package main

import (
	"context"
	"time"
		"fmt"

		"gopkg.in/olivere/elastic.v5"
	"encoding/json"
	"reflect"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`




func main() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	url:="http://hadoop:9200/"
	client:=getClient(url)

	Ping(url,client,ctx)

	createIndex(client,ctx)


	//Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	setDataJson("twitter","tweet","1",tweet1,client,ctx)


	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	setData("twitter","tweet","2",tweet2,client,ctx)



	getData("twitter","tweet","1",client,ctx)





	// Flush to make sure the documents got written.
	client.Flush().Index("twitter").Do(ctx)


	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)


	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index
			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}
			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}


	script:=elastic.NewScriptInline("ctx._source.retweets += params.num").
		Lang("painless").
		Param("num", 1)
	doc:=map[string]interface{}{"retweets": 0}
	update("twitter","tweet","1",script,doc,client,ctx)

	deleteIndex(client,ctx)

}

func getClient(url string) *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		// Must turn off sniff in docker
		elastic.SetSniff(false),
		//elastic.SetBasicAuth("user", "secret"),		//用户名密码访问
	)
	if err != nil {
		// Handle error
		panic(err)
	}
	return client
}

func Ping(url string,client *elastic.Client,ctx context.Context)  {
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	esversion, err := client.ElasticsearchVersion(url)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

}

func createIndex(client *elastic.Client,ctx context.Context)  {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("twitter").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}
func deleteIndex(client *elastic.Client,ctx context.Context)  {
	// Delete an index.
	deleteIndex, err := client.DeleteIndex("twitter").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}
func setDataJson(indexName,typeName,id string,body interface{} ,client *elastic.Client,ctx context.Context){
	put1, err := client.Index().
		Index(indexName).
		Type(typeName).
		Id(id).
		BodyJson(body).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func setData(indexName,typeName,id,body string ,client *elastic.Client,ctx context.Context){
	put1, err := client.Index().
		Index(indexName).
		Type(typeName).
		Id(id).
		BodyString(body).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}


func getData(indexName,typeName,id string ,client *elastic.Client,ctx context.Context)  {
	// Get tweet with specified ID
	get1, err := client.Get().
		Index(indexName).
		Type(typeName).
		Id(id).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

}



func update(indexName,typeName,id string ,script *elastic.Script,doc interface{},client *elastic.Client,ctx context.Context)  {
	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	update, err := client.Update().Index(indexName).Type(typeName).Id(id).
		Script(script).
		Upsert(doc).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("New version of tweet %q is now %d\n", update.Id, update.Version)
}