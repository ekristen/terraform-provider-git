package provider

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccGitRepositoryDataSourceConfigBasic(path string) string {
	return fmt.Sprintf(`
data "git_repository" "test" {
  path = %[1]q
}
`, path)
}

func testAccGitRepositoryDataSourceConfigComplex(path string) string {
	return fmt.Sprintf(`
data "git_repository" "test" {
  path = %[1]q
  ref_short_length = 9
}
`, path)
}

func TestAccGitRepositoryDataSource1(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "terraform-provider-git-")
	assert.NoError(t, err)
	//noinspection GoUnhandledErrorResult
	defer os.RemoveAll(tempDir)

	hash, err := testSetupGit(tempDir, "", 0)
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigBasic(tempDir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", tempDir),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "has_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "ref", hash.String()),
					resource.TestCheckResourceAttr("data.git_repository.test", "ref_short", hash.String()[0:7]),
				),
			},
		},
	})
}

func TestAccGitRepositoryDataSource2(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "terraform-provider-git-")
	assert.NoError(t, err)
	//noinspection GoUnhandledErrorResult
	defer os.RemoveAll(tempDir)

	hash, err := testSetupGit(tempDir, "", 0)
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigComplex(tempDir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", tempDir),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "has_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "ref", hash.String()),
					resource.TestCheckResourceAttr("data.git_repository.test", "ref_short", hash.String()[0:9]),
				),
			},
		},
	})
}

func TestAccGitRepositoryDataSource3(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "terraform-provider-git-")
	assert.NoError(t, err)
	//noinspection GoUnhandledErrorResult
	defer os.RemoveAll(tempDir)

	hash, err := testSetupGit(tempDir, "v1.0.0", 1)
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigBasic(tempDir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", tempDir),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "has_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "semver", fmt.Sprintf("v1.0.0-1.g%s", hash.String()[0:7])),
					resource.TestCheckResourceAttr("data.git_repository.test", "ref", hash.String()),
				),
			},
		},
	})
}

func TestAccGitRepositoryDataSource4(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "terraform-provider-git-")
	assert.NoError(t, err)
	//noinspection GoUnhandledErrorResult
	defer os.RemoveAll(tempDir)

	hash, err := testSetupGit(tempDir, "v1.0.0", 0)
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigBasic(tempDir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", tempDir),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "has_tag", "true"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "semver", "v1.0.0"),
					resource.TestCheckResourceAttr("data.git_repository.test", "ref", hash.String()),
				),
			},
		},
	})
}

func TestAccGitRepositoryDataSource5(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "terraform-provider-git-")
	assert.NoError(t, err)
	//noinspection GoUnhandledErrorResult
	defer os.RemoveAll(tempDir)

	reg, err := regexp.Compile("unable to open git repository")
	assert.NoError(t, err)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccGitRepositoryDataSourceConfigBasic(tempDir),
				ExpectError: reg,
			},
		},
	})
}

func testSetupGit(path string, tag string, extraCommits int) (*plumbing.Hash, error) {
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return nil, err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(path, "README.md"), []byte("testing"), 0644); err != nil {
		return nil, err
	}

	if _, err := wt.Add("README.md"); err != nil {
		return nil, err
	}

	hash, err := wt.Commit("tests", &git.CommitOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	if tag != "" {
		ref, err := repo.CreateTag(tag, hash, &git.CreateTagOptions{
			Message: tag,
		})
		if err != nil {
			return nil, err
		}
		_ = ref
	}

	for i := 0; i < extraCommits; i++ {
		if err := os.WriteFile(filepath.Join(path, "README.md"), []byte(fmt.Sprintf("testing %02d", i)), 0644); err != nil {
			return nil, err
		}
		if _, err := wt.Add("README.md"); err != nil {
			return nil, err
		}
		hash, err = wt.Commit(fmt.Sprintf("tests %02d", i), &git.CommitOptions{
			All: true,
		})
		if err != nil {
			return nil, err
		}
	}

	return &hash, nil
}
