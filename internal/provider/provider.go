package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure CronMathProvider satisfies various provider interfaces.
var _ provider.Provider = &CronMathProvider{}

// CronMathProvider defines the provider implementation.
type CronMathProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CronMathProviderModel describes the provider data model.
type CronMathProviderModel struct{}

func (p *CronMathProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cronmath"
	resp.Version = p.version
}

func (p *CronMathProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:          map[string]schema.Attribute{},
		MarkdownDescription: "The cronmath provider allows manipulation of cron expressions using Terraform.",
	}
}

func (p *CronMathProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CronMathProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *CronMathProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCronScheduleResource,
	}
}

func (p *CronMathProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCronMathDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CronMathProvider{
			version: version,
		}
	}
}
