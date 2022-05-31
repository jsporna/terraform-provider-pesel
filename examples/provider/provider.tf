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
