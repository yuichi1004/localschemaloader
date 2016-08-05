package localschemaloader

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonreference"
	"github.com/yuichi1004/gojsonschema"
)

type LocalSchemaLoaderFactory struct {
	gojsonschema.JSONLoaderFactory
	urlBasePath  string
	fileBasePath string
}

type LocalSchemaLoader struct {
	gojsonschema.JSONLoader
	source  string
	factory LocalSchemaLoaderFactory
}

func New(urlBasePath, fileBasePath string) LocalSchemaLoaderFactory {
	return LocalSchemaLoaderFactory{
		urlBasePath:  urlBasePath,
		fileBasePath: fileBasePath,
	}
}

func (f LocalSchemaLoaderFactory) New(source string) gojsonschema.JSONLoader {
	if f.urlBasePath == "" {
		panic("")
	}
	return &LocalSchemaLoader{
		source:  source,
		factory: f,
	}
}

func (l *LocalSchemaLoader) JsonSource() interface{} {
	return l.source
}
func (l *LocalSchemaLoader) LoadJSON() (interface{}, error) {
	reference, err := l.JsonReference()
	if err != nil {
		return nil, err
	}

	refToUrl := reference
	refToUrl.GetUrl().Fragment = ""
	source := refToUrl.String()
	dest := strings.Replace(source, l.factory.urlBasePath, l.factory.fileBasePath, 1)

	return l.loadFromFile(dest)
}

func (l *LocalSchemaLoader) JsonReference() (gojsonreference.JsonReference, error) {
	return gojsonreference.NewJsonReference(l.source)
}

func (l *LocalSchemaLoader) LoaderFactory() gojsonschema.JSONLoaderFactory {
	return l.factory
}

func (l *LocalSchemaLoader) loadFromFile(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return decodeJsonUsingNumber(f)
}

func decodeJsonUsingNumber(r io.Reader) (interface{}, error) {

	var document interface{}

	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}

	return document, nil

}
