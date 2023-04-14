# one2sentinel

A Go program that exports 1Password usage, signin and audit events to Microsoft Sentinel SIEM.

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
go run ./cmd/... -config=dev.yml
INFO[0000] shipping logs                                 module=sentinel_logs table_name=OnePasswordLogs total=82
INFO[0002] shipped logs                                  module=sentinel_logs table_name=OnePasswordLogs
INFO[0002] successfully sent logs to sentinel            total=82
```

Or binary:
```shell
% one2sen -config=config.yml
```

## Building

```shell
% make build
```
