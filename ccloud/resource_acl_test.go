package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_BasicACL(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(testACL_noConfig, u),
			},
			{
				ResourceName:      "confluentcloud_acl.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//lintignore:AT004
const testACL_noConfig = `
resource "confluentcloud_acl" "test" {
	cluster_id  		= "lkc-v9ky0"
	bootstrap_servers 	= "https://pkac-57298.eu-west-1.aws.confluent.cloud"
	resource_type		= "GROUP"
	pattern_type 		= "LITERAL"
	name        		= "acl-test-%s"
	principal			= "User:1522"
	operation			= "READ"
	host				= "*"
	permission_type		= "ALLOW"
}
`
