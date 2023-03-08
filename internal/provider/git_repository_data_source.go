package provider

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	gitutils "github.com/ekristen/terraform-provider-git/pkg/git"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GitRepository{}

func NewGitRepository() datasource.DataSource {
	return &GitRepository{}
}

// GitRepository defines the data source implementation.
type GitRepository struct {
	client *http.Client
}

// GitRepositoryModel describes the data source data model.
type GitRepositoryModel struct {
	Path              types.String `tfsdk:"path"`
	Branch            types.String `tfsdk:"branch"`
	Summary           types.String `tfsdk:"summary"`
	Semver            types.String `tfsdk:"semver"`
	SemverFallbackTag types.String `tfsdk:"semver_fallback_tag"`
}

func (d *GitRepository) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *GitRepository) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Git Repository data source",

		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Path to Git Repository",
				Required:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch Name",
				Computed:            true,
			},
			"semver": schema.StringAttribute{
				MarkdownDescription: "Git Summary in SEMVER format",
				Computed:            true,
			},
			"summary": schema.StringAttribute{
				MarkdownDescription: "Git Summary",
				Computed:            true,
			},
			"semver_fallback_tag": schema.StringAttribute{
				MarkdownDescription: "Fallback Tag for SEMVER Generation",
				Optional:            true,
			},
		},
	}
}

func (d *GitRepository) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *GitRepository) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GitRepositoryModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.SemverFallbackTag.ValueString() == "" {
		data.SemverFallbackTag = types.StringValue("v0.0.0")
	}

	repo, err := git.PlainOpen(data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Git Error", err.Error())
		return
	}

	head, err := repo.Head()
	if err != nil {
		resp.Diagnostics.AddError("Git Error", err.Error())
		return
	}

	tagName, counter, headHash, err := gitutils.Describe(*repo)

	result, err := gitutils.GenerateVersion(*tagName, *counter, *headHash, time.Now(), gitutils.GenerateVersionOptions{
		FallbackTagName: data.SemverFallbackTag.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Generate Version Error", err.Error())
		return
	}

	worktree, _ := repo.Worktree()
	status, _ := worktree.Status()
	dirty := !status.IsClean()

	if tagName != nil && toString(tagName) != "" {
		data.Summary = types.StringValue(fmt.Sprintf("%s-%d-g%s", toString(tagName), toInt(counter), toString(headHash)[0:7]))
	} else {
		data.Summary = types.StringValue(fmt.Sprintf("%s", toString(headHash)[0:7]))
	}

	if dirty {
		data.Summary = types.StringValue(fmt.Sprintf("%s-dirty", data.Summary.ValueString()))
	}

	data.Semver = types.StringValue(*result)
	data.Branch = types.StringValue(head.Name().String())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func toString(original *string) string {
	if original != nil {
		return *original
	}
	return ""
}

func toInt(original *int) int {
	if original != nil {
		return *original
	}
	return 0
}
