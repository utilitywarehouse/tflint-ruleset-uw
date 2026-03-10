# TFLint UW ruleset

This is a tflint plugin for UW specific needs.

In most cases, teams can run terraform commands on their local machines, so
linting is not a security enforcer, but a nicer UI for the rules enforced by
other security methods (like permission boundaries).

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

| Name | Description |
| --- | --- |
| [`aws_s3_bucket_policy`](rules/aws_s3_bucket_policy.md) | Enforces UW S3 bucket naming and required `Name` tag |
| [`aws_security_group_policy`](rules/aws_security_group_policy.md) | Requires aws security groups to have an "owner" tag |

## Using the plugin locally

Clone the repository locally and run the following command:

```
$ go test ./...
$ make install
```

Add the local (no version/source attributes) version of the plugin to your .tflint.hcl

```
plugin "uw" {
  enabled = true
}
```

## Releasing

Releases are created via github. Creating a new release on github will trigger
a workflow that uses goreleaser to create the required builds and checksums.
