package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccExampleDataSourceConfig1 = `
data "git_repository" "test" {
  path = "./testdata/no-tags"
}
`

func TestAccGitRepositoryDataSource1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", "./testdata/no-tags"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig2 = `
data "git_repository" "test" {
  path = "./testdata/tagged"
}
`

func TestAccGitRepositoryDataSource2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", "./testdata/tagged"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig3 = `
data "git_repository" "test" {
  path = "./testdata/tagged-extra-commits"
}
`

func TestAccGitRepositoryDataSource3(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "id", "./testdata/tagged-extra-commits"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_tag", "false"),
					resource.TestCheckResourceAttr("data.git_repository.test", "is_dirty", "false"),
				),
			},
		},
	})
}
