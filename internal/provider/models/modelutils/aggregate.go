package modelutils

var Aggregate_Operations = map[string]string{
	"add":                        "DEFAULT",
	"sum":                        "SUM",
	"minimum":                    "MIN",
	"maximum":                    "MAX",
	"average":                    "AVG",
	"set_intersection":           "SET_INTERSECTION",
	"distribution_concatenation": "DIST_CONCAT",
	"custom":                     "CUSTOM",
}
