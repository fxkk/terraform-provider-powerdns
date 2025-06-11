---
layout: "powerdns"
page_title: "PowerDNS: powerdns_zone_metadata"
sidebar_current: "docs-powerdns-zone-metadata"
description: |-
  Manages PowerDNS zone metadata.
---

# powerdns\_zone\_metadata

Provides a PowerDNS zone metadata resource.

## Example Usage

```hcl
# Add metadata to a zone
resource "powerdns_zone" "example" {
  name        = "example.com."
  kind        = "Native"
  nameservers = ["ns1.example.com.", "ns2.example.com."]
}

resource "powerdns_zone_metadata" "test" {

   metadata {
    kind = "ALLOW-AXFR-FROM"
    values = ["AUTO-NS"]
  }
  
}

```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The name of the zone to which the metadata belongs.
* `metadata` - (Required) A `metadata` block as documented below.

---

A `metadata` block supports the following arguments:
* `kind` - (Required) The kind of metadata.
* `values` - (Required) A list of values for the metadata kind.



## Importing

An existing zone metadata can be imported into this resource by supplying the zone name and metadata kind.

When importing this resource, only the metadata kinds explicitly defined in your Terraform configuration will be managed. Other existing metadata for the zone that is not declared in your configuration will remain untouched.

When first importing this resource, Terraform will attempt to apply changes to the metadata kinds specified in your configuration. This happens because the provider in this framework version cannot read fields directly from the resource config before it exists in the state.

```
$ terraform import powerdns_zone_metadata.test example.com.
```

For more information on how to use terraform's `import` command, please refer to terraform's [core documentation](https://www.terraform.io/docs/import/index.html#currently-state-only).
