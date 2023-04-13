# one2sentinel

A Go program that exports 1Password usage and signin events to Microsoft Sentinel SIEM.

## Running

First create a yaml file, such as `config.yml`:
```yaml
log:
  level: INFO

microsoft:
  app_id: ""
  secret_key: ""
  tenant_id: ""
  subscription_id: ""
  resource_group: ""
  workspace_name: ""
  workspace_id: ""
  workspace_primary_key: ""
  expires_months: 6
  update_table: false

onepassword:
  api_token: ""
```

And now run the program from source code:
```shell
% make
```

Or binary:
```shell
% one2sen -config=config.yml
```

## Building

```shell
% make build
```
