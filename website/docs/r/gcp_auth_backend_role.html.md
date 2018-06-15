---
layout: "vault"
page_title: "Vault: vault_auth_backend resource"
sidebar_current: "docs-vault-resource-gcp-auth-backend-role"
description: |-
  Managing roles in an GCP auth backend in Vault
---

# vault\_gcp\_auth\_backend\_role

Provides a resource to create a role in an [GCP auth backend within Vault](https://www.vaultproject.io/docs/auth/gcp.html).

## Example Usage

```hcl
resource "vault_auth_backend" "gcp" {
    path = "gcp"
    type = "gcp"
}

resource "vault_gcp_auth_backend_role" "gcp" {
    backend                = "${vault_auth_backend.cert.path}"
    project_id             = "foo-bar-baz"
    bound_service_accounts = ["database-server@foo-bar-baz.iam.gserviceaccount.com"]
    policies               = ["database-server"]

}
```

## Argument Reference

The following arguments are supported:

* `role` - (Required) Name of the GCP role

* `type` - (Required) Type of GCP authentication role

* `project_id` - (Required) GCP Project that the role exists within

* `ttl` - (Optional) Default TTL of tokens issued by the backend

* `max_ttl` - (Optional) Maximum TTL of tokens issued by the backend

* `period` - (Optional) Duration in seconds for token.  If set, the issued token is a periodic token.

* `policies` - (Optional) Policies to grant on the issued token

* `bound_service_accounts` - (Optional) GCP Service Accounts allowed to issue tokens under this role

* `backend` - (Optional) Path to the mounted GCP auth backend

For more details on the usage of each argument consult the [Vault GCP API documentation](https://www.vaultproject.io/api/auth/gcp/index.html).

## Attribute Reference

No additional attributes are exposed by this resource.
