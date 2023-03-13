package provider

import (
	"fmt"
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

func TestAccGitRepositoryDataSource1(t *testing.T) {
	path := "./testdata/no-tags"
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigBasic(path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", path),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
				),
			},
		},
	})
}

func TestAccGitRepositoryDataSource2(t *testing.T) {
	path := "./testdata/tagged"
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigBasic(path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", path),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "semver", "v1.0.0"),
				),
			},
		},
	})
}

func TestAccGitRepositoryDataSource3(t *testing.T) {
	path := "./testdata/tagged-extra-commits"
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGitRepositoryDataSourceConfigBasic(path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", path),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "semver", "v1.0.0-1.g21e4385"),
				),
			},
		},
	})
}
