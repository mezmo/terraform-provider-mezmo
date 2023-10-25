package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/mezmo/terraform-provider-mezmo/internal/provider"
)

var (
	_, b, _, _       = runtime.Caller(0)
	basepath         = filepath.Dir(b)
	testdataPath     = path.Join(basepath, "testdata")
	processorsPath   = path.Join(testdataPath, "processors")
	sourcesPath      = path.Join(testdataPath, "sources")
	destinationsPath = path.Join(testdataPath, "destinations")
)

func loadJsonFile[T any](t *testing.T, baseDir string, filename string) *T {
	t.Helper()
	bytes, err := os.ReadFile(path.Join(baseDir, filename))
	if err != nil {
		t.Fatalf("Could not read %s file. Reason: %s", filename, err)
	}
	var into T
	if err := json.Unmarshal(bytes, &into); err != nil {
		t.Fatalf("Could not encode %s to json. Reason: %s", filename, err)
	}
	return &into
}

func TestConvertToTerraformModel(t *testing.T) {
	convResources, err := ConvertibleResources()

	if err != nil {
		t.Fatalf("error retrieving convertible resources. reason: %s", err)
	}
	resList := loadResources(t)

	for _, res := range convResources {
		t.Run(fmt.Sprintf("convert %s resource", res.TypeName()), func(t *testing.T) {
			m, ok := resList[res.TypeName()]
			if !ok {
				t.Errorf("resource %s not found", res.TypeName())
			}
			_, err := res.ConvertToTerraformModel(m)
			if err != nil {
				t.Errorf("failed to convert %s to terraform model. reason: %s", res.TypeName(), err)
			}
		})
	}
}

func loadResources(t *testing.T) map[string]*reflect.Value {
	resList := make(map[string]*reflect.Value)
	pipeline := reflect.ValueOf(*loadJsonFile[PipelineApiModel](t, testdataPath, "pipeline.json"))
	resList["pipeline"] = &pipeline
	addToMap(t, resList, loadDirFiles[ProcessorApiModel](t, processorsPath, func(filename string) string {
		return fmt.Sprintf("%s_processor", filename)
	}))
	addToMap(t, resList, loadDirFiles[SourceApiModel](t, sourcesPath, func(filename string) string {
		return fmt.Sprintf("%s_source", filename)
	}))
	addToMap(t, resList, loadDirFiles[DestinationApiModel](t, destinationsPath, func(filename string) string {
		return fmt.Sprintf("%s_destination", filename)
	}))
	return resList
}

func loadDirFiles[T any](t *testing.T, dirPath string, mapKeyFn func(filename string) string) map[string]*reflect.Value {
	dirFiles := make(map[string]*reflect.Value)
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fExt := filepath.Ext(info.Name())
		if fExt != ".json" {
			return nil
		}
		mVal := reflect.ValueOf(*loadJsonFile[T](t, dirPath, info.Name()))
		fBaseName := strings.Split(info.Name(), ".json")[0]
		dirFiles[mapKeyFn(fBaseName)] = &mVal
		return nil
	})
	if err != nil {
		t.Fatalf("could not read files from %s. reason: %s", dirPath, err)
	}
	return dirFiles
}

func addToMap[K string, T any](t *testing.T, target map[K]*T, source map[K]*T) {
	for k, v := range source {
		target[k] = v
	}
}

func TestConvertibleResources(t *testing.T) {
	var unmatchedResources = make(map[string]bool)

	p := provider.MezmoProvider{}
	for _, resFn := range p.Resources(context.Background()) {
		res := resFn()
		resTy := reflect.TypeOf(res).String()
		unmatchedResources[resTy] = false
	}

	actual, err := ConvertibleResources()
	if err != nil {
		t.Errorf("Error calling ConvertibleResources: %s", err)
		return
	}
	for _, convRes := range actual {
		convResTy := reflect.TypeOf(convRes).String()
		if _, ok := unmatchedResources[convResTy]; ok {
			delete(unmatchedResources, convResTy)
		} else {
			unmatchedResources[convResTy] = true
		}
	}

	for unmatchedType, found := range unmatchedResources {
		if !found {
			t.Errorf(
				"Terraform Resource %s does not contain corresponding Convertible Resource instance",
				unmatchedType,
			)
		} else {
			t.Errorf(
				"Convertible Resource %s does not contain corresponding Terraform Resource",
				unmatchedType,
			)
		}
	}
}

func findConvertibleResource(t *testing.T, resourceTypeName string) ConvertibleResourceDef {
	t.Helper()
	convResources, err := ConvertibleResources()
	if err != nil {
		t.Fatal("no convertible resources found")
	}

	for _, res := range convResources {
		if res.TypeName() == resourceTypeName {
			return res
		}
	}
	t.Fatalf("resource %s not found", resourceTypeName)
	return nil
}
