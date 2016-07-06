package localschemaloader

import (
	"log"
	"testing"

	"github.com/yuichi1004/gojsonschema"
)

type TestObj struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}

func TestBasic(t *testing.T) {
	factory := New(
		"https://github.com/yuichi1004/localschemaloader/",
		"./",
	)

	test := TestObj{
		Id:    1,
		Value: "tiny",
	}

	schemaLoader := factory.New("https://github.com/yuichi1004/localschemaloader/test/scheme1.json")
	documentLoader := gojsonschema.NewGoLoader(test)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		log.Printf("The document is valid\n")
	} else {
		log.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			t.Errorf("- %s\n", desc)
		}
	}

}
