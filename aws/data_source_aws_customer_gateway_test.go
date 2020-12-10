package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSCustomerGatewayDataSource_Filter(t *testing.T) {
	dataSourceName := "data.aws_customer_gateway.test"
	resourceName := "aws_customer_gateway.test"

	asn := acctest.RandIntRange(64512, 65534)
	hostOctet := acctest.RandIntRange(1, 254)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCustomerGatewayDataSourceConfigFilter(asn, hostOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bgp_asn", dataSourceName, "bgp_asn"),
					resource.TestCheckResourceAttrPair(resourceName, "ip_address", dataSourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.%", dataSourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
				),
			},
		},
	})
}

func TestAccAWSCustomerGatewayDataSource_ID(t *testing.T) {
	dataSourceName := "data.aws_customer_gateway.test"
	resourceName := "aws_customer_gateway.test"

	asn := acctest.RandIntRange(64512, 65534)
	hostOctet := acctest.RandIntRange(1, 254)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCustomerGatewayDataSourceConfigID(asn, hostOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bgp_asn", dataSourceName, "bgp_asn"),
					resource.TestCheckResourceAttrPair(resourceName, "ip_address", dataSourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.%", dataSourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
				),
			},
		},
	})
}

func TestAccAWSCustomerGatewayDataSource_CertificateArn(t *testing.T) {
	dataSourceName := "data.aws_customer_gateway.test"
	resourceName := "aws_customer_gateway.test"

	asn := acctest.RandIntRange(64512, 65534)
	hostOctet := acctest.RandIntRange(1, 254)
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCustomerGatewayDataSourceConfigCertificateArn(rName, asn, hostOctet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bgp_asn", dataSourceName, "bgp_asn"),
					resource.TestCheckResourceAttrPair(resourceName, "ip_address", dataSourceName, "ip_address"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.%", dataSourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_arn", dataSourceName, "certificate_arn"),
				),
			},
		},
	})
}

func testAccAWSCustomerGatewayDataSourceConfigFilter(asn, hostOctet int) string {
	name := acctest.RandomWithPrefix("test-filter")
	return fmt.Sprintf(`
resource "aws_customer_gateway" "test" {
  bgp_asn    = %d
  ip_address = "50.0.0.%d"
  type       = "ipsec.1"

  tags = {
    Name = "%s"
  }
}

data "aws_customer_gateway" "test" {
  filter {
    name   = "tag:Name"
    values = [aws_customer_gateway.test.tags.Name]
  }
}
`, asn, hostOctet, name)
}

func testAccAWSCustomerGatewayDataSourceConfigID(asn, hostOctet int) string {
	return fmt.Sprintf(`
resource "aws_customer_gateway" "test" {
  bgp_asn    = %d
  ip_address = "50.0.0.%d"
  type       = "ipsec.1"
}

data "aws_customer_gateway" "test" {
  id = aws_customer_gateway.test.id
}
`, asn, hostOctet)
}

func testAccAWSCustomerGatewayDataSourceConfigCertificateArn(rName string, asn, hostOctet int) string {
	return fmt.Sprintf(`
resource "aws_acmpca_certificate_authority" "test" {
	permanent_deletion_time_in_days = 7
	type                            = "ROOT"

	certificate_authority_configuration {
		key_algorithm     = "RSA_4096"
		signing_algorithm = "SHA512WITHRSA"

		subject {
			common_name = "terraformtesting.com"
		}
	}
}

resource "aws_acm_certificate" "cert" {
	domain_name               = "%s.terraformtesting.com"
	certificate_authority_arn = aws_acmpca_certificate_authority.test.arn
}

resource "aws_customer_gateway" "test" {
  bgp_asn         = %d
  ip_address      = "50.0.0.%d"
  certificate_arn = aws_acm_certificate.test.arn
  type            = "ipsec.1"
}

data "aws_customer_gateway" "test" {
  id = aws_customer_gateway.test.id
}
`, rName, asn, hostOctet)
}
