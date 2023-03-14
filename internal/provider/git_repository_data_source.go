package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	Id                   types.String `tfsdk:"id"`
	Path                 types.String `tfsdk:"path"`
	Reference            types.String `tfsdk:"ref"`
	ReferenceShort       types.String `tfsdk:"ref_short"`
	Summary              types.String `tfsdk:"summary"`
	Branch               types.String `tfsdk:"branch"`
	Tag                  types.String `tfsdk:"tag"`
	IsDirty              types.Bool   `tfsdk:"is_dirty"`
	IsTag                types.Bool   `tfsdk:"is_tag"`
	IsBranch             types.Bool   `tfsdk:"is_branch"`
	HasTag               types.Bool   `tfsdk:"has_tag"`
	CommitCount          types.Int64  `tfsdk:"commit_count"`
	Semver               types.String `tfsdk:"semver"`
	SemverFallbackTag    types.String `tfsdk:"semver_fallback_tag"`
	ReferenceShortLength types.Int64  `tfsdk:"ref_short_length"`
}

func (d *GitRepository) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *GitRepository) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Git Repository data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "id",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path to Git Repository",
				Required:            true,
			},
			"summary": schema.StringAttribute{
				MarkdownDescription: "Git Summary",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch Name",
				Computed:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Current Tag of Repository",
				Computed:            true,
			},
			"ref": schema.StringAttribute{
				MarkdownDescription: "Current reference of the repository",
				Computed:            true,
			},
			"ref_short": schema.StringAttribute{
				MarkdownDescription: "Short version of the current reference",
				Computed:            true,
			},
			"ref_short_length": schema.Int64Attribute{
				MarkdownDescription: "Length of the short version of the current reference (default: 7)",
				Optional:            true,
			},
			"is_branch": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the current reference is a branch",
				Computed:            true,
			},
			"is_dirty": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the repository is in a dirty state",
				Computed:            true,
			},
			"is_tag": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the current reference is a tag",
				Computed:            true,
			},
			"has_tag": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the current reference has been tagged",
				Computed:            true,
			},
			"commit_count": schema.Int64Attribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"semver": schema.StringAttribute{
				MarkdownDescription: "Git Summary in SEMVER format",
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
	if data.ReferenceShortLength.ValueInt64() == 0 {
		data.ReferenceShortLength = types.Int64Value(7)
	}

	repo, err := git.PlainOpen(data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to open git repository", err.Error())
		return
	}

	head, err := repo.Head()
	if err != nil {
		resp.Diagnostics.AddError("unable to read git head reference", err.Error())
		return
	}

	tagName, counter, headHash, err := gitutils.Describe(*repo)
	if err != nil {
		resp.Diagnostics.AddError("unable to run git describe", err.Error())
		return
	}

	data.Reference = types.StringValue(head.Hash().String())
	data.ReferenceShort = types.StringValue(head.Hash().String()[0:data.ReferenceShortLength.ValueInt64()])
	data.CommitCount = types.Int64Value(int64(*counter))

	result, err := gitutils.GenerateVersion(*tagName, *counter, *headHash, time.Now(), gitutils.GenerateVersionOptions{
		FallbackTagName: data.SemverFallbackTag.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("unable to generate version", err.Error())
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		resp.Diagnostics.AddError("unable to read worktree", err.Error())
		return
	}

	status, err := worktree.Status()
	if err != nil {
		resp.Diagnostics.AddError("unable to get worktree status", err.Error())
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("head: %s", head.Name().String()))
	tflog.Trace(ctx, fmt.Sprintf("is_tag: %t", head.Name().IsTag()))
	tflog.Trace(ctx, fmt.Sprintf("is_branch: %t", head.Name().IsBranch()))

	dirty := !status.IsClean()
	isTag := head.Name().IsTag()

	if tagName != nil && toString(tagName) != "" {
		data.Summary = types.StringValue(fmt.Sprintf("%s-%d-g%s", toString(tagName), toInt(counter), toString(headHash)[0:7]))
	} else {
		data.Summary = types.StringValue(fmt.Sprintf("%s", toString(headHash)[0:7]))
	}

	if dirty {
		data.Summary = types.StringValue(fmt.Sprintf("%s-dirty", data.Summary.ValueString()))
	}

	data.HasTag = types.BoolValue(false) // default
	iter, err := repo.Tags()
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		if ref == nil {
			return nil
		}

		tflog.Trace(ctx, fmt.Sprintf("ref: %s", ref.Hash().String()))

		obj, err := repo.TagObject(ref.Hash())
		if err != nil && !errors.Is(err, plumbing.ErrObjectNotFound) {
			return err
		}

		if obj == nil {
			return nil
		}

		if obj.Target.String() == head.Hash().String() {
			data.HasTag = types.BoolValue(true)
		}
		return nil
	}); err != nil {
		resp.Diagnostics.AddError("unable to find tag for reference", err.Error())
		return
	}

	data.Id = types.StringValue(data.Path.ValueString())
	data.Semver = types.StringValue(*result)
	data.Branch = types.StringValue(head.Name().String())
	data.IsDirty = types.BoolValue(dirty)
	data.IsTag = types.BoolValue(isTag)
	data.IsBranch = types.BoolValue(head.Name().IsBranch())

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
