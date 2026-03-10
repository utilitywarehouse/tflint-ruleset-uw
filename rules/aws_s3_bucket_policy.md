# aws_s3_bucket_policy

S3 buckets in `dev` and `prod` must follow UW naming conventions and include a `Name` tag.

Note that the rule accepts team prefixes as valid, but the errors only suggest full team names, as that's the preferred style.

## Why
- Bucket name defines ownership and access
- `Name` is the chosen tag for detailed cost allocation on S3
