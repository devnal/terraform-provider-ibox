package ibox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPDNSRecord_A(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigA,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists("powerdns_record.test-a"),
				),
			},
		},
	})
}

func TestAccPDNSRecord_WithCount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigHyphenedWithCount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists("powerdns_record.test-counted.0"),
					testAccCheckPDNSRecordExists("powerdns_record.test-counted.1"),
				),
			},
		},
	})
}

func testAccCheckPDNSRecordDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns_record" {
			continue
		}

		client := testAccProvider.Meta().(*Client)
		exists, err := client.RecordExistsByID(rs.Primary.Attributes["zone"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error checking if record still exists: %#v", rs.Primary.ID)
		}
		if exists {
			return fmt.Errorf("Record still exists: %#v", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckPDNSRecordExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		foundRecords, err := client.ListRecordsByID(rs.Primary.Attributes["zone"], rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(foundRecords) == 0 {
			return fmt.Errorf("Record does not exist")
		}
		for _, rec := range foundRecords {
			if rec.Id() == rs.Primary.ID {
				return nil
			}
		}
		return fmt.Errorf("Record does not exist: %#v", rs.Primary.ID)
	}
}

const testPDNSRecordConfigA = `
resource "powerdns_record" "test-a" {
  zone = "sysa.xyz"
    name = "redis.sysa.xyz"
    type = "A"
    ttl = 60
    records = [ "1.1.1.1", "2.2.2.2" ]
}`

const testPDNSRecordConfigHyphenedWithCount = `
resource "powerdns_record" "test-counted" {
    count = "2"
    zone = "sysa.xyz"
    name = "redis-${count.index}.sysa.xyz"
    type = "A"
    ttl = 60
    records = [ "1.1.1.${count.index}" ]
}`
