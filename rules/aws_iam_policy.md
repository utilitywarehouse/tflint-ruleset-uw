# aws_iam_policy

In `dev` and `prod`, iam resources must follow specific naming patterns and have a permission boundary set.

Ensures IAM users, roles, groups, policies, and instance profiles use `name`
values that start with `$team-` or `$team_prefix-`.

## Why
Naming conventions and permission boundaries are how we implement access control in AWS
