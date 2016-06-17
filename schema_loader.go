package localschemaloder

import (
	"bytes"
	"os"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
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
	source string
	factory LocalSchemaLoaderFactory
}

func New(urlBasePath, fileBasePath string) LocalSchemaLoaderFactory {
	return LocalSchemaLoaderFactory{
		urlBasePath: urlBasePath,
		fileBasePath: fileBasePath,
	}
}


func (f LocalSchemaLoaderFactory) New(source string) gojsonschema.JSONLoader {
	if f.urlBasePath == "" {
		panic("")
	}
	return &LocalSchemaLoader{
		source: source,
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
	document, err := l.loadFromFile(dest)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func (l *LocalSchemaLoader) JsonReference() (gojsonreference.JsonReference, error) {
	return gojsonreference.NewJsonReference(l.JsonSource().(string))
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

	bodyBuff, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return decodeJsonUsingNumber(bytes.NewReader(bodyBuff))

}

func decodeJsonUsingNumber(r io.Reader) (interface{}, error) {

	var document interface{}

	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	err := decoder.Decode(&document)
	if err != nil {
		return nil, err
	}

	return document, nil

}

