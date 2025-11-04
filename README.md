# Grist API library

Use Grist API natively with Go.

[⚠️WIP] [🛑Not ready to use]

## Getting started

* Have a running Grist instance 
`$ docker compose up -d -f build/docker-compose.yml`

Run examples locally 
`$ GRIST_API_KEY=<API_KEY_FROM_GRIST> mage orgs`
`$ GRIST_API_KEY=<API_KEY_FROM_GRIST> mage worskpaces`

* `GRIST_API_KEY` must be generated directly from Grist settings on the WebUI

TODO: 
* Orgs 🛠️
  * List ✅
  * Describe ✅
  * Modify ✅
  * Delete ✅
  * List users access ⚠️ (API does not match documentation, open PR if needed)
  * Edit users access 🛑
* Workspaces 
    * List ✅
    * Describe ✅
    * Modify ✅
    * Delete ✅
* Docs
    * Describe ✅
    * ModifyMetadata ✅
    * Delete ✅
* Records 🛑
* Tables 🛑
* Columns 🛑
* Attachments 🛑
* Webhooks 🛑
* SQL 🛑
* Users 🛑
* SCIM 🛑

