package resources

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/processors"
)

var (
	_, b, _, _   = runtime.Caller(0)
	basepath     = filepath.Dir(b)
	testdataPath = path.Join(basepath, "testdata")
)

func loadJsonFile[T any](t *testing.T, filename string) *T {
	t.Helper()
	bytes, err := os.ReadFile(path.Join(testdataPath, filename))
	if err != nil {
		t.Fatalf("Could not read %s file. Reason: %s", filename, err)
	}
	var into T
	if err := json.Unmarshal(bytes, &into); err != nil {
		t.Fatalf("Could not encode %s to json. Reason: %s", filename, err)
	}
	return &into
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

	for unmatchedType, which := range unmatchedResources {
		if which == false {
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

func TestRouteResource(t *testing.T) {
	type args struct {
		resourceName string
		model        *reflect.Value
	}
	tests := []struct {
		name       string
		args       args
		assertions func(p *reflect.Value)
	}{
		{
			name: "converts simple route processor to terraform model",
			args: args{
				resourceName: "route_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "simple_route.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.RouteProcessorModel)

				var actualInputs []string
				for _, item := range model.Inputs.Elements() {
					v, _ := item.(basetypes.StringValue)
					actualInputs = append(actualInputs, v.ValueString())
				}
				expectedInputs := []string{"12d94f50-6c68-11ee-bdee-6671faf7df66"}
				if !reflect.DeepEqual(actualInputs, expectedInputs) {
					t.Errorf("simple_route processor inputs do not match.\nexpected: %s\ngot: %s", expectedInputs, actualInputs)
				}

				if len(model.Conditionals.Elements()) == 0 {
					t.Fatal("simple_route processor returned zero conditionals. wanted 1")
				}

				elem, _ := model.Conditionals.Elements()[0].(basetypes.ObjectValue)
				label := elem.Attributes()["label"].(basetypes.StringValue).ValueString()
				logicalOperation := elem.Attributes()["logical_operation"].(basetypes.StringValue).ValueString()
				outputName := elem.Attributes()["output_name"].(basetypes.StringValue).ValueString()

				if label != "Error logs" {
					t.Errorf("simple_route processor label error. wanted: 'Error logs', got: '%s'", label)
				}
				if logicalOperation != "OR" {
					t.Errorf("simple_route processor logical operation error. wanted: 'OR', got: '%s'", logicalOperation)
				}
				if outputName != "805821a7" {
					t.Errorf("simple_route processor output name error. wanted: '805821a7', got: '%s'", outputName)
				}
			},
		},
		{
			name: "converts route processor to terraform model",
			args: args{
				resourceName: "route_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "route.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.RouteProcessorModel)

				var actualInputs []string
				for _, item := range model.Inputs.Elements() {
					v, _ := item.(basetypes.StringValue)
					actualInputs = append(actualInputs, v.ValueString())
				}
				expectedInputs := []string{
					"12d94f50-6c68-11ee-bdee-6671faf7df66",
					"7b212506-23cb-11ed-b300-4ef12c27e273",
				}
				if !reflect.DeepEqual(actualInputs, expectedInputs) {
					t.Errorf("route processor inputs do not match.\nexpected: %s\ngot: %s", expectedInputs, actualInputs)
				}

				if len(model.Conditionals.Elements()) == 0 {
					t.Fatal("route processor returned zero conditionals. wanted 2")
				}

				expectedConditionals := []struct {
					label            string
					outputName       string
					logicalOperation string
				}{
					{
						label:            "Error logs",
						outputName:       "805821a7",
						logicalOperation: "OR",
					},
					{
						label:            "App info logs",
						outputName:       "c6b6ebe5",
						logicalOperation: "AND",
					},
				}
				for i, v := range model.Conditionals.Elements() {
					elem, _ := v.(basetypes.ObjectValue)
					label := elem.Attributes()["label"].(basetypes.StringValue).ValueString()
					logicalOperation := elem.Attributes()["logical_operation"].(basetypes.StringValue).ValueString()
					outputName := elem.Attributes()["output_name"].(basetypes.StringValue).ValueString()

					if label != expectedConditionals[i].label {
						t.Errorf("route processor label error. wanted: '%s', got: '%s'", expectedConditionals[i].label, label)
					}
					if logicalOperation != expectedConditionals[i].logicalOperation {
						t.Errorf("route processor logical operation error. wanted: '%s', got: '%s'", expectedConditionals[i].logicalOperation, logicalOperation)
					}
					if outputName != expectedConditionals[i].outputName {
						t.Errorf("route processor logical operation error. wanted: '%s', got: '%s'", expectedConditionals[i].logicalOperation, logicalOperation)
					}
				}
			},
		},
		{
			name: "converts parse sequentially processor to terraform model",
			args: args{
				resourceName: "parse_sequentially_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "parse_sequentially.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.ParseSequentiallyProcessorModel)

				var actualInputs []string
				for _, item := range model.Inputs.Elements() {
					v, _ := item.(basetypes.StringValue)
					actualInputs = append(actualInputs, v.ValueString())
				}
				expectedInputs := []string{
					"110ec7da-5c7e-11ee-bffb-26dab184329f",
				}
				if !reflect.DeepEqual(actualInputs, expectedInputs) {
					t.Errorf("parse sequentially inputs do not match.\nexpected: %s\ngot: %s", expectedInputs, actualInputs)
				}

				if len(model.Parsers.Elements()) == 0 {
					t.Fatal("parse sequentially returned zero parsers. wanted 2")
				}

				expectedParsers := []struct {
					label      string
					parser     string
					outputName string
				}{
					{
						label:      "Apache Error",
						outputName: "36d9714483c9745012cd14f9380335ac",
						// converts from parse_apache_log to apache_log
						parser: "apache_log",
					},
					{
						label:      "Nginx Combined",
						outputName: "5db52e644356529da4e34663969833b9",
						// converts from parse_nginx_log to nginx_log
						parser: "nginx_log",
					},
				}
				for i, v := range model.Parsers.Elements() {
					elem, _ := v.(basetypes.ObjectValue)
					label := elem.Attributes()["label"].(basetypes.StringValue).ValueString()
					parser := elem.Attributes()["parser"].(basetypes.StringValue).ValueString()
					outputName := elem.Attributes()["output_name"].(basetypes.StringValue).ValueString()

					if label != expectedParsers[i].label {
						t.Errorf("parse sequentially label error. wanted: '%s', got: '%s'", expectedParsers[i].label, label)
					}
					if parser != expectedParsers[i].parser {
						t.Errorf("parse sequentially parser error. wanted: '%s', got: '%s'", expectedParsers[i].parser, parser)
					}
					if outputName != expectedParsers[i].outputName {
						t.Errorf("parse sequentially output name error. wanted: '%s', got: '%s'", expectedParsers[i].outputName, outputName)
					}
				}
			},
		},
		{
			name: "converts dedupe processor to terraform model",
			args: args{
				resourceName: "dedupe_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "dedupe.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.DedupeProcessorModel)

				expected := struct {
					numEvents      int64
					comparisonType string
					inputs         []string
					fields         []string
				}{
					numEvents:      5000,
					comparisonType: "Match",
					inputs:         []string{"7b212506-23cb-11ed-b300-4ef12c27e273"},
					fields:         []string{".foo", ".bar"},
				}

				numEvents := model.NumberOfEvents.ValueInt64()
				if numEvents != expected.numEvents {
					t.Errorf("dedupe number of events don't match. wanted: %v, got: %v", expected.numEvents, numEvents)
				}
				compType := model.ComparisonType.ValueString()
				if compType != expected.comparisonType {
					t.Errorf("dedupe comparison types don't match. wanted: %s, got: %s", expected.comparisonType, compType)
				}
				var fields []string
				for _, item := range model.Fields.Elements() {
					fields = append(fields, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(fields, expected.fields) {
					t.Errorf("dedupe fields don't match. wanted: %s, got: %s", expected.fields, fields)
				}
				var inputs []string
				for _, item := range model.Inputs.Elements() {
					inputs = append(inputs, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(inputs, expected.inputs) {
					t.Errorf("dedupe inputs don't match. wanted: %s, got: %s", expected.inputs, inputs)
				}
			},
		},
		{
			name: "converts drop fields processor to terraform model",
			args: args{
				resourceName: "drop_fields_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "drop_fields.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.DropFieldsProcessorModel)

				expected := struct {
					inputs []string
					fields []string
				}{
					inputs: []string{"7b212506-23cb-11ed-b300-4ef12c27e273"},
					fields: []string{".errors", ".warnings"},
				}

				var fields []string
				for _, item := range model.Fields.Elements() {
					fields = append(fields, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(fields, expected.fields) {
					t.Errorf("drop processor fields don't match. wanted: %s, got: %s", expected.fields, fields)
				}
				var inputs []string
				for _, item := range model.Inputs.Elements() {
					inputs = append(inputs, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(inputs, expected.inputs) {
					t.Errorf("drop processor inputs don't match. wanted: %s, got: %s", expected.inputs, inputs)
				}
			},
		},
		{
			name: "converts flatten fields processor to terraform model",
			args: args{
				resourceName: "flatten_fields_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "flatten_fields.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.FlattenFieldsProcessorModel)

				expected := struct {
					inputs    []string
					fields    []string
					delimiter string
				}{
					inputs:    []string{"7b212506-23cb-11ed-b300-4ef12c27e273"},
					fields:    []string{".list", ".map"},
					delimiter: ",",
				}

				delimiter := model.Delimiter.ValueString()
				if delimiter != expected.delimiter {
					t.Errorf("flatten fields delimiters don't match. wanted: %s, got: %s", expected.delimiter, delimiter)
				}

				var fields []string
				for _, item := range model.Fields.Elements() {
					fields = append(fields, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(fields, expected.fields) {
					t.Errorf("flatten fields processor fields don't match. wanted: %s, got: %s", expected.fields, fields)
				}
				var inputs []string
				for _, item := range model.Inputs.Elements() {
					inputs = append(inputs, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(inputs, expected.inputs) {
					t.Errorf("flatten fields processor inputs don't match. wanted: %s, got: %s", expected.inputs, inputs)
				}
			},
		},
		{
			name: "converts map fields processor to terraform model",
			args: args{
				resourceName: "map_fields_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "map_fields.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.MapFieldsProcessorModel)

				type mapping struct {
					target          string
					source          string
					dropSource      bool
					overwriteTarget bool
				}
				expected := struct {
					inputs   []string
					mappings []mapping
				}{
					inputs: []string{"7b212506-23cb-11ed-b300-4ef12c27e273"},
					mappings: []mapping{
						{
							dropSource:      true,
							source:          ".firstname",
							target:          ".fname",
							overwriteTarget: false,
						},
						{
							dropSource:      true,
							source:          ".lastname",
							target:          ".lname",
							overwriteTarget: true,
						},
					},
				}

				var inputs []string
				for _, item := range model.Inputs.Elements() {
					inputs = append(inputs, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(inputs, expected.inputs) {
					t.Errorf("map fields processor inputs don't match. wanted: %s, got: %s", expected.inputs, inputs)
				}

				for i, item := range model.Mappings.Elements() {
					v := item.(basetypes.ObjectValue)
					source := v.Attributes()["source_field"].(basetypes.StringValue).ValueString()
					target := v.Attributes()["target_field"].(basetypes.StringValue).ValueString()
					dropSource := v.Attributes()["drop_source"].(basetypes.BoolValue).ValueBool()
					overwriteTarget := v.Attributes()["overwrite_target"].(basetypes.BoolValue).ValueBool()

					if source != expected.mappings[i].source {
						t.Errorf("map fields sources don't match. wanted: %s, got: %s", expected.mappings[i].source, source)
					}
					if target != expected.mappings[i].target {
						t.Errorf("map fields targets don't match. wanted: %s, got: %s", expected.mappings[i].target, target)
					}
					if dropSource != expected.mappings[i].dropSource {
						t.Errorf("map fields drop sources don't match. wanted: %v, got: %v", expected.mappings[i].dropSource, dropSource)
					}
					if overwriteTarget != expected.mappings[i].overwriteTarget {
						t.Errorf("map fields overwrite targets don't match. wanted: %v, got: %v", expected.mappings[i].overwriteTarget, overwriteTarget)
					}
				}
			},
		},
		{
			name: "converts parse processor to terraform model",
			args: args{
				resourceName: "parse_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "parse.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.ParseProcessorModel)

				expected := struct {
					inputs  []string
					field   string
					target  string
					parser  string
					options map[string]string
				}{
					inputs: []string{"7b212506-23cb-11ed-b300-4ef12c27e273"},
					field:  ".",
					target: ".parsed",
					parser: "nginx_log",
					options: map[string]string{
						"format":           "error",
						"timestamp_format": "%Y-%m-%dT%H:%M:%SZ",
					},
				}

				var inputs []string
				for _, item := range model.Inputs.Elements() {
					inputs = append(inputs, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(inputs, expected.inputs) {
					t.Errorf("parse processor inputs don't match. wanted: %s, got: %s", expected.inputs, inputs)
				}

				if model.Field.ValueString() != expected.field {
					t.Errorf("parse processor fields don't match. wanted: %s, got: %s", expected.field, model.Field.ValueString())
				}
				if model.TargetField.ValueString() != expected.target {
					t.Errorf("parse processor target fields don't match. wanted: %s, got: %s", expected.target, model.TargetField.ValueString())
				}
				format := model.NginxOptions.Attributes()["format"].(basetypes.StringValue).ValueString()
				timeFormat := model.NginxOptions.Attributes()["timestamp_format"].(basetypes.StringValue).ValueString()

				if format != expected.options["format"] {
					t.Errorf("parse processor format options don't match. wanted: %s, got: %s", expected.options["format"], format)
				}
				if timeFormat != expected.options["timestamp_format"] {
					t.Errorf("parse processor timestamp format options don't match. wanted: %s, got: %s", expected.options["timestamp_format"], timeFormat)
				}
			},
		},
		{
			name: "converts reduce processor to terraform model",
			args: args{
				resourceName: "reduce_processor",
				model: func() *reflect.Value {
					t.Helper()
					model := loadJsonFile[ProcessorApiModel](t, "reduce.json")
					rModel := reflect.ValueOf(*model)
					return &rModel
				}(),
			},
			assertions: func(p *reflect.Value) {
				model := p.Interface().(processors.ReduceProcessorModel)

				type dateFormat struct {
					field  string
					format string
				}
				type mergeStrategy struct {
					field    string
					strategy string
				}
				expected := struct {
					inputs          []string
					groupBy         []string
					dateFormats     []dateFormat
					mergeStrategies []mergeStrategy
				}{
					inputs: []string{"7b212506-23cb-11ed-b300-4ef12c27e273"},
					groupBy: []string{
						".error.level",
						".user.email",
					},
					dateFormats: []dateFormat{
						{
							field:  ".log_date",
							format: "%Y-%m-%dT%H:%M:%S",
						},
					},
					mergeStrategies: []mergeStrategy{
						{
							field:    ".errors",
							strategy: "flat_unique",
						},
						{
							field:    ".users",
							strategy: "concat_raw",
						},
					},
				}

				var inputs []string
				for _, item := range model.Inputs.Elements() {
					inputs = append(inputs, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(inputs, expected.inputs) {
					t.Errorf("reduce processor inputs don't match. wanted: %s, got: %s", expected.inputs, inputs)
				}

				var dateFormats []dateFormat
				for _, item := range model.DateFormats.Elements() {
					v := item.(basetypes.ObjectValue)
					dateFormats = append(dateFormats, dateFormat{
						field:  v.Attributes()["field"].(basetypes.StringValue).ValueString(),
						format: v.Attributes()["format"].(basetypes.StringValue).ValueString(),
					})
				}
				if !reflect.DeepEqual(dateFormats, expected.dateFormats) {
					t.Errorf("reduce processor date formats don't match. wanted: %s, got: %s", expected.dateFormats, dateFormats)
				}

				var mergeStrategies []mergeStrategy
				for _, item := range model.MergeStrategies.Elements() {
					v := item.(basetypes.ObjectValue)
					mergeStrategies = append(mergeStrategies, mergeStrategy{
						field:    v.Attributes()["field"].(basetypes.StringValue).ValueString(),
						strategy: v.Attributes()["strategy"].(basetypes.StringValue).ValueString(),
					})
				}
				if !reflect.DeepEqual(mergeStrategies, expected.mergeStrategies) {
					t.Errorf("reduce processor date formats don't match. wanted: %s, got: %s", expected.mergeStrategies, mergeStrategies)
				}

				var groupBy []string
				for _, item := range model.GroupBy.Elements() {
					groupBy = append(groupBy, item.(basetypes.StringValue).ValueString())
				}
				if !reflect.DeepEqual(groupBy, expected.groupBy) {
					t.Errorf("reduce processor group bys don't match. wanted: %s, got: %s", expected.groupBy, groupBy)
				}

				conditions := model.FlushCondition.Attributes()
				flushWhen := conditions["when"].(basetypes.StringValue).ValueString()
				if flushWhen != "starts_when" {
					t.Errorf("reduce processor flush whens don't match. wanted: starts_when, got: %s", flushWhen)
				}
				// 2 expressions and 1 nested expression
				expressions := conditions["conditional"].(basetypes.ObjectValue).Attributes()["expressions"].(basetypes.ListValue)
				expressionGroups := conditions["conditional"].(basetypes.ObjectValue).Attributes()["expressions_group"].(basetypes.ListValue)

				if len(expressionGroups.Elements()) != 1 {
					t.Errorf("reduce processor returned %v nested expressions. wanted: 1", len(expressionGroups.Elements()))
				}
				if len(expressions.Elements()) != 2 {
					t.Errorf("reduce processor returned %v expressions. wanted: 2", len(expressions.Elements()))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := findConvertibleResource(t, tt.args.resourceName)
			m, err := res.ConvertToTerraformModel(tt.args.model)

			if err != nil {
				t.Fatalf("could not convert %s to terraform model. reason: %s", tt.args.resourceName, err)
			}
			tt.assertions(m)
		})
	}
}
