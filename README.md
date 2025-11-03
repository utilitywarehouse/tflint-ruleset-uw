# TFLint UW ruleset

This is a tflint plugin for enforing UW rules

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "uw" {
  enabled = true

  version = "x.y.z"
  source  = "github.com/utilitywarehouse/tflint-ruleset-uw"
}
```

## Rules

| Name                                                                                    | Description                                         |
|-----------------------------------------------------------------------------------------|-----------------------------------------------------|
| [`aws_security_group_missing_owner_tag`](rules/aws_security_group_missing_owner_tag.md) | Requires aws security groups to have an "owner" tag |

## Using the plugin locally

Clone the repository locally and run the following command:

```
$ make install
```

Add the local (no "source" attribute) version of the plugin to your .tflint.hcl

```
plugin "uw" {
  enabled = true
}
```
