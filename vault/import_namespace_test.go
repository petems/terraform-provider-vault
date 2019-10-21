package vault

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNamespace_importBasic(t *testing.T) {

	isEnterprise := os.Getenv("TF_ACC_ENTERPRISE")
	if isEnterprise == "" {
		t.Skip("TF_ACC_ENTERPRISE is not set, test is applicable only for Enterprise version of Vault")
	}

	path := "test-" + acctest.RandString(10)
	cfg := namespaceConfig{
		path: path,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testNamespaceConfig(cfg.path),
				Check:  testNestedNamespaceCheckAttrs(cfg.path),
			},
			{
				ResourceName:      "vault_namespace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
