package alicloud

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsDomainsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_alidns_domains.default", strconv.FormatInt(int64(rand), 10), dataSourceAlidnsDomainsConfigDependence)
	aliDomainConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "false",
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "true",
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dns_domain.default.id}-fake"},
		}),
	}
	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "false",
			"group_name_regex":  "${alicloud_dns_group.default.name}",
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ali_domain":        "false",
			"group_name_regex":  "${alicloud_dns_group.default.name}_fake",
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"instance_id":       "",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"instance_id":       "fake",
		}),
	}
	versionCodeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"version_code":      "mianfei",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"version_code":      "bumianfei",
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"resource_group_id": fmt.Sprintf("%s_fake", os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"ids":               []string{"${alicloud_dns_domain.default.id}"},
			"version_code":      "mianfei",
			"instance_id":       "",
			"ali_domain":        "false",
			"group_name_regex":  "${alicloud_dns_group.default.name}",
			"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name_regex": "${alicloud_dns_domain.default.domain_name}",
			"ids":               []string{"${alicloud_dns_domain.default.id}"},
			"version_code":      "mianfei",
			"instance_id":       "",
			"ali_domain":        "true",
			"group_name_regex":  "${alicloud_dns_group.default.name}",
			"resource_group_id": os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
		}),
	}
	AlidnsDomainsCheckInfo.dataSourceTestCheck(t, rand, aliDomainConf, idsConf, groupNameConf, instanceIdConf, versionCodeConf, resourceGroupIdConf, allConf)
}

func dataSourceAlidnsDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "dnsName"{
	default = "tf-testacc%sdnsalidomainbasic%s.abc"
}

variable "dnsGroupName"{
	default = "tf-testaccdns%s"
}

resource "alicloud_dns_group" "default" {
  name = "${var.dnsGroupName}"
}

resource "alicloud_dns_domain" "default" {
	domain_name = "${var.dnsName}"
	group_id = "${alicloud_dns_group.default.id}"
	resource_group_id = "%s"
}
`, defaultRegionToTest, name, name, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

var existAlidnsDomainsMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"domains.#":               "1",
		"domains.0.domain_id":     CHECKSET,
		"domains.0.domain_name":   fmt.Sprintf("tf-testacc%sdnsalidomainbasic%d.abc", defaultRegionToTest, rand),
		"domains.0.ali_domain":    "false",
		"domains.0.group_id":      CHECKSET,
		"domains.0.group_name":    fmt.Sprintf("tf-testaccdns%d", rand),
		"domains.0.instance_id":   "",
		"domains.0.version_code":  "mianfei",
		"domains.0.puny_code":     CHECKSET,
		"domains.0.dns_servers.#": CHECKSET,
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"names.#":                 "1",
		"names.0":                 fmt.Sprintf("tf-testacc%sdnsalidomainbasic%d.abc", defaultRegionToTest, rand),
	}
}

var fakeAlidnsDomainsMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"names.#":   "0",
		"ids.#":     "0",
		"domains.#": "0",
	}
}

var AlidnsDomainsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_alidns_domains.default",
	existMapFunc: existAlidnsDomainsMapCheck,
	fakeMapFunc:  fakeAlidnsDomainsMapCheck,
}
