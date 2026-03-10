# aws_security_group_policy

Security groups in `dev` and `prod` must have an "owner" tag with one of the UW teams

## Why
Existing security groups can only be edited by the owner team, so this rules
ensures that teams won't create security groups that they can't manage later
