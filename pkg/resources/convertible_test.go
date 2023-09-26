package resources

import (
	"context"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider"
	"reflect"
	"testing"
)

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
