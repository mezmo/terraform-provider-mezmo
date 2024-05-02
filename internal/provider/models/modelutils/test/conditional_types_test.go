package modelutils_test

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	"github.com/stretchr/testify/assert"
)

var (
	_, b, _, _   = runtime.Caller(0)
	basepath     = filepath.Dir(b)
	testCaseFile = path.Join(basepath, "..", "testdata", "unwind_conditional.json")
)

func TestUnwindConditionalToModel(t *testing.T) {
	testCases, err := LoadTestCases()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range testCases {
		t.Run(tt.Description, func(t *testing.T) {
			res := modelutils.UnwindConditionalToModel(tt.UserConfig["conditional"].(map[string]any), modelutils.Non_Change_Operator_Labels)
			// Terraform attr.Value String method returns null as <null> and unknown as <unknown>
			actual_string := strings.ReplaceAll(res.String(), "<null>", "\"null\"")
			actual_string = strings.ReplaceAll(actual_string, "<unknown>", "\"unknown\"")

			var actual any
			json.Unmarshal([]byte(actual_string), &actual)

			assert.EqualValues(t, tt.ExpectedExpr, actual)
		})
	}
}

type testCase struct {
	Description  string         `json:"description"`
	UserConfig   map[string]any `json:"user_config"`
	ExpectedExpr any            `json:"expected_expression"`
}

func LoadTestCases() ([]testCase, error) {
	contents, err := os.ReadFile(testCaseFile)
	if err != nil {
		return nil, err
	}
	var testCases []testCase
	json.Unmarshal(contents, &testCases)
	return testCases, nil
}
