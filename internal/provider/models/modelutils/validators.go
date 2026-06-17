package modelutils

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type uriValidator struct{}

func (v uriValidator) Description(_ context.Context) string {
	return "value must be a valid http or https URL"
}

func (v uriValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uriValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			"Attribute "+req.Path.String()+" must be a valid http or https URL, got: "+value,
		)
	}
}

// URIValidator validates that a string attribute is an absolute http or https URL.
func URIValidator() validator.String {
	return uriValidator{}
}
