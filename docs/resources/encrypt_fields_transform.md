---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mezmo_encrypt_fields_transform Resource - terraform-provider-mezmo"
subcategory: ""
description: |-
  Encrypts the value of the provided field
---

# mezmo_encrypt_fields_transform (Resource)

Encrypts the value of the provided field



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `algorithm` (String) The encryption algorithm to use on the field
- `field` (String) Field to encrypt. The value of the field must be a primitive (string, number, boolean).
- `iv_field` (String) The field in which to store the generated initialization vector, IV. Each encrypted value will have a unique IV.
- `key` (String, Sensitive) The encryption key
- `pipeline_id` (String) The uuid of the pipeline

### Optional

- `description` (String) A user-defined value describing the transform component
- `encode_raw_bytes` (Boolean) Encode the encrypted value and generated initialization vector as Base64 text
- `inputs` (List of String) The ids of the input components
- `title` (String) A user-defined title for the transform component

### Read-Only

- `generation_id` (Number) An internal field used for component versioning
- `id` (String) The uuid of the transform component

