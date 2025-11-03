# aws_security_group_missing_owner_tag

Security groups must have an "owner" tag

## Why
Existing security groups can only be edited by the owner team, so this rules
ensures that teams won't create security groups that they can't manage later
