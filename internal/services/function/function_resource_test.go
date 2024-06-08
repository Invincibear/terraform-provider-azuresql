package function_test

import (
	"fmt"
	"terraform-provider-azuresql/internal/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type FunctionResource struct{}

func TestAccCreateFunctionBasic(t *testing.T) {
	acceptance.PreCheck(t)
	data := acceptance.BuildTestData(t)
	r := FunctionResource{}

	connections := []string{
		data.FabricDatabase_connection,
		data.SQLDatabase_connection,
		data.SynapseDatabase_connection,
	}

	for _, connection := range connections {
		print(fmt.Sprintf("\n\nRunning test for connection %s\n\n", connection))
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:                   r.basic(connection, data.RandomString),
					ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
				},
				{
					Config:                   r.basic(connection, data.RandomString),
					ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
					ResourceName:             "azuresql_function.test",
					ImportState:              true,
					ImportStateVerify:        true,
				},
			},
		})
	}
}

func TestAccCreateFunctionProps(t *testing.T) {
	acceptance.PreCheck(t)
	data := acceptance.BuildTestData(t)
	r := FunctionResource{}

	connections := []string{
		data.SQLDatabase_connection,
	}

	for _, connection := range connections {
		print(fmt.Sprintf("\n\nRunning test for connection %s\n\n", connection))
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:                   r.propsapi(connection, data.RandomString),
					ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
				},
			},
		})
	}
}

func (r FunctionResource) basic(connection string, name string) string {
	template := r.template()

	return fmt.Sprintf(
		`
		%[1]s

		data "azuresql_schema" "dbo" {
			database 	= "%[2]s"
			name 		= "dbo"
		}

		resource "azuresql_function" "test" {
			database 	= "%[2]s"
			name        = "tffunction_%[3]s"
			schema		= data.azuresql_schema.dbo.id
			raw         = <<-EOT
				create FUNCTION dbo.tffunction_%[3]s()
				returns table 
				with SCHEMABINDING AS
				return  
				select 1 as a
			EOT
		}
		`, template, connection, name)
}

func (r FunctionResource) propsapi(connection string, name string) string {
	template := r.template()

	return fmt.Sprintf(
		`
		%[1]s

		data "azuresql_schema" "dbo" {
			database 	= "%[2]s"
			name 		= "dbo"
		}

		resource "azuresql_function" "test" {
			database 	= "%[2]s"
			name        = "tffunction_%[3]s"
			schema		= data.azuresql_schema.dbo.id
			properties = {
				arguments = [
				  {
					name = "a"
					type = "int"
				  },
				  {
					name = "b"
					type = "int"
				  }
				]
				executor      = "self"
				return_type   = "int"
				schemabinding = true
				definition    = "@a + @b"
			  }
		}
		`, template, connection, name)
}

func (r FunctionResource) template() string {
	return fmt.Sprintf(`
		provider "azuresql" {
		}
	`)
}
