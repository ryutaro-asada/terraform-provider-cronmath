package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ryutaro-asada/cronmath"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CronScheduleResource{}
var _ resource.ResourceWithImportState = &CronScheduleResource{}

func NewCronScheduleResource() resource.Resource {
	return &CronScheduleResource{}
}

// CronScheduleResource defines the resource implementation.
type CronScheduleResource struct{}

// CronScheduleResourceModel describes the resource data model.
type CronScheduleResourceModel struct {
	Id          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	BaseCron    types.String      `tfsdk:"base_cron"`
	FinalCron   types.String      `tfsdk:"final_cron"`
	Description types.String      `tfsdk:"description"`
	Adjustments []AdjustmentModel `tfsdk:"adjustments"`
}

type AdjustmentModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.Int64  `tfsdk:"value"`
	Unit  types.String `tfsdk:"unit"`
}

func (r *CronScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *CronScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Manages a cron schedule with time adjustments",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the schedule",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the cron schedule",
				Required:            true,
			},
			"base_cron": schema.StringAttribute{
				MarkdownDescription: "The base cron expression to start from",
				Required:            true,
			},
			"final_cron": schema.StringAttribute{
				MarkdownDescription: "The final cron expression after all adjustments",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of what this schedule is for",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"adjustments": schema.ListNestedBlock{
				MarkdownDescription: "List of time adjustments to apply to the base cron",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of adjustment: 'add' or 'sub'",
							Required:            true,
						},
						"value": schema.Int64Attribute{
							MarkdownDescription: "The value to adjust by",
							Required:            true,
						},
						"unit": schema.StringAttribute{
							MarkdownDescription: "Unit of time: 'minutes' or 'hours'",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *CronScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (r *CronScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CronScheduleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Creating cron schedule", map[string]interface{}{
		"name":      data.Name.ValueString(),
		"base_cron": data.BaseCron.ValueString(),
	})

	// Calculate final cron
	finalCron, err := r.calculateFinalCron(ctx, data.BaseCron.ValueString(), data.Adjustments)
	if err != nil {
		resp.Diagnostics.AddError("Calculation Error", err.Error())
		return
	}

	// Generate ID
	data.Id = types.StringValue(fmt.Sprintf("cron_%s", data.Name.ValueString()))
	data.FinalCron = types.StringValue(finalCron)

	tflog.Trace(ctx, "Created cron schedule", map[string]interface{}{
		"id":         data.Id.ValueString(),
		"final_cron": data.FinalCron.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CronScheduleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Recalculate final cron to ensure consistency
	finalCron, err := r.calculateFinalCron(ctx, data.BaseCron.ValueString(), data.Adjustments)
	if err != nil {
		resp.Diagnostics.AddError("Calculation Error", err.Error())
		return
	}

	data.FinalCron = types.StringValue(finalCron)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CronScheduleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Updating cron schedule", map[string]interface{}{
		"id":        data.Id.ValueString(),
		"base_cron": data.BaseCron.ValueString(),
	})

	// Recalculate final cron
	finalCron, err := r.calculateFinalCron(ctx, data.BaseCron.ValueString(), data.Adjustments)
	if err != nil {
		resp.Diagnostics.AddError("Calculation Error", err.Error())
		return
	}

	data.FinalCron = types.StringValue(finalCron)

	tflog.Trace(ctx, "Updated cron schedule", map[string]interface{}{
		"final_cron": data.FinalCron.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CronScheduleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Deleted cron schedule", map[string]interface{}{
		"id": data.Id.ValueString(),
	})
}

func (r *CronScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CronScheduleResource) calculateFinalCron(ctx context.Context, baseCron string, adjustments []AdjustmentModel) (string, error) {
	// Use the external cronmath package
	cm := cronmath.New(baseCron)

	for i, adj := range adjustments {
		value := int(adj.Value.ValueInt64())
		unit := adj.Unit.ValueString()
		adjType := adj.Type.ValueString()

		tflog.Trace(ctx, fmt.Sprintf("Applying adjustment %d", i), map[string]interface{}{
			"type":  adjType,
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
			return "", fmt.Errorf("invalid unit: %s", unit)
		}

		switch adjType {
		case "add":
			cm = cm.Add(duration)
		case "sub":
			cm = cm.Sub(duration)
		default:
			return "", fmt.Errorf("invalid adjustment type: %s", adjType)
		}
	}

	if err := cm.Error(); err != nil {
		return "", fmt.Errorf("failed to calculate cron: %w", err)
	}

	return cm.String(), nil
}
