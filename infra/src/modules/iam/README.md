# Module for Aws IAM resources

Module to create AWS IAM related resources to be implemented for Enron's mail indexing solution on ZincSearch

## Available Features

- IAM role

## Usage

### Generate task execution role

```hcl
module "ecs" {
  source = "/modules/iam"

  # required
  task_execution_role_name = "ecsTaskExecutionRole"
  
  # optionals
  policy_attachment_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  region_name = "us-east-1"
}
```

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3.7 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_policy_document.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | resource |
| [aws_iam_role.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_task_execution_role_name"></a> [task_execution_role_name](#input_task_execution_role_name) | Role name for ecs task execution | `string` | `""` | yes |
| <a name="input_policy_attachment_arn"></a> [policy_attachment_arn](#input_policy_attachment_arn) | Policy attachment arn for task execution role | `string` | `"arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"` | no |
| <a name="input_region_name"></a> [region_name](#input_region_name) | Aws region name | `string` | `"us-east-1"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_iam_ecs_task_execution_role_arn"></a> [iam_ecs_task_execution_role_arn](#outp_iam_ecs_task_execution_role_arn) | Arn of the role for execution of ecs tasks |

## Authors

Module is maintained by [WitsoftGroup](https://github.com/WitsoftGroup)
