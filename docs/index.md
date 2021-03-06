---
page_title: "Provider: PESEL"
description: |-
  The PESEL provider is used to generate PESEL identification number.
---

# PESEL Provider

The "pesel" provider allows the use of PESEL id number within Terraform
configurations. This is a *logical provider*, which means that it works
entirely within Terraform's logic, and doesn't interact with any other
services.

To force a pesel result to be replaced, the `taint` command can be used to
produce a new result on the next run.

For example:
```terraform
terraform {
  required_providers {
    pesel = {
      source = "jsporna/pesel"
    }
  }
}

resource "pesel_id" "random" {
}

data "pesel_id" "somebody" {
  id = "65432101239"
}

output "random" {
  value = pesel_id.random
}

output "somebody" {
  value = data.pesel_id.somebody
}
```