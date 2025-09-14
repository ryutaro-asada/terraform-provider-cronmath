package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ryutaro-asada/cronmath"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CronMathDataSource{}

func NewCronMathDataSource() datasource.DataSource {
	return &CronMathDataSource{}
}

// CronMathDataSource defines the data source implementation.
type CronMathDataSource struct{}

// CronMathDataSourceModel describes the data source data model.
type CronMathDataSourceModel struct {
	Id         types.String     `tfsdk:"id"`
	Input      types.String     `tfsdk:"input"`
	Operations []OperationModel `tfsdk:"operations"`
	Result     types.String     `tfsdk:"result"`
}

type OperationModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (d *CronMathDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_calculate"
}

func (d *CronMathDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Calculate cron expressions with time arithmetic operations",

		Attributes: map[string]schema.Attribute{
			"input": schema.StringAttribute{
				MarkdownDescription: "The input cron expression (5 fields: minute hour day month weekday)",
				Required:            true,
			},
			"result": schema.StringAttribute{
				MarkdownDescription: "The resulting cron expression after applying all operations",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this calculation",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"operations": schema.ListNestedBlock{
				MarkdownDescription: "List of operations to apply to the cron expression",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Operation type: 'add' or 'sub'",
							Required:            true,
						},
						"value": schema.Int64Attribute{
							MarkdownDescription: "The value to add or subtract",
							Required:            true,
						},
						"unit": schema.StringAttribute{
							MarkdownDescription: "Time unit: 'minutes' or 'hours'",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (d *CronMathDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *CronMathDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CronMathDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Processing cron expression", map[string]interface{}{
		"input":            data.Input.ValueString(),
		"operations_count": len(data.Operations),
	})

	// Create CronMath instance using the external package
	cm := cronmath.New(data.Input.ValueString())

	// Apply each operation
	for i, op := range data.Operations {
		value := int(op.Value.ValueInt64())
		unit := op.Unit.ValueString()
		opType := op.Type.ValueString()

		tflog.Trace(ctx, fmt.Sprintf("Applying operation %d", i), map[string]interface{}{
			"type":  opType,
			"value": value,
			"unit":  unit,
		})

		var duration cronmath.Duration
		switch unit {
		case "minutes", "minute", "min", "m":
			duration = cronmath.Minutes(value)
		case "hours", "hour", "hr", "h":
			duration = cronmath.Hours(value)
		default:
			resp.Diagnostics.AddError(
				"Invalid Unit",
				fmt.Sprintf("Unit must be 'minutes' or 'hours', got '%s'", unit),
			)
			return
		}

		switch opType {
		case "add":
			cm = cm.Add(duration)
		case "sub":
			cm = cm.Sub(duration)
		default:
			resp.Diagnostics.AddError(
				"Invalid Operation",
				fmt.Sprintf("Operation type must be 'add' or 'sub', got '%s'", opType),
			)
			return
		}
	}

	// Check for errors
	if err := cm.Error(); err != nil {
		resp.Diagnostics.AddError(
			"Calculation Error",
			fmt.Sprintf("Failed to calculate cron expression: %s", err),
		)
		return
	}

	// Set the result
	data.Result = types.StringValue(cm.String())
	data.Id = types.StringValue(fmt.Sprintf("%s-%d", data.Input.ValueString(), len(data.Operations)))

	tflog.Trace(ctx, "Cron calculation complete", map[string]interface{}{
		"result": data.Result.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
