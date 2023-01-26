# main.tf

### local definitions ###
locals {
  # iam
  task_role_name   = "indexEcsTaskExecutionRole-${var.ENV}"
  lambda_role_name = "indexLambdaExecutionRole-${var.ENV}"

  # nerwork
  lb_name      = "zinc-lb-${var.ENV}"
  zinc_tg_name = "zinc-tg-${var.ENV}"

  # ecs
  cluster_name  = "zinc-cluster-${var.ENV}"
  task_name     = "zinc-task-${var.ENV}"
  task_img_name = "public.ecr.aws/zinclabs/zinc:latest"
  service_name  = "zinc-service-${var.ENV}"

  # s3
  data_bucket_name = "enron-mail-data-${var.ENV}"

  # Lambda
  lambda_name    = "indez-eron-email-${var.ENV}"
  lambda_runtime = "go1.x"
  lambda_handler = "main"
}
#########################

### data ###
data "archive_file" "dummy_lambda" {
  type        = "zip"
  output_path = "${path.module}/lambda_function_payload.zip"

  source {
    content  = "hello lambda!"
    filename = "dummy_lambda.txt"
  }
}
############

### local modules ###
module "iam" {
  # setup
  source = "./modules/iam"
  for_each = tomap({
    "${local.task_role_name}" = {
      identifiers  = ["ecs-tasks.amazonaws.com"]
      arn_policies = ["arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
    },
    "${local.lambda_role_name}" = {
      identifiers = ["lambda.amazonaws.com"]
      arn_policies = [
        "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
        "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
      ]
    }
  })

  # variables
  region_name  = var.AWS_REGION
  role_name    = each.key
  identifiers  = each.value.identifiers
  arn_policies = each.value.arn_policies
}

module "network" {
  # setup
  source = "./modules/network"
  depends_on = [
    module.iam
  ]

  # variables
  region_name       = var.AWS_REGION
  lb_name           = local.lb_name
  target_group_name = local.zinc_tg_name
}

module "ecs" {
  # setup
  source = "./modules/ecs"
  depends_on = [
    module.network
  ]

  # variables
  cluster_name            = local.cluster_name
  task_name               = local.task_name
  task_img_name           = local.task_img_name
  target_group_arn        = module.network.lb_target_group_arn
  role_task_execution_arn = module.iam["${local.task_role_name}"].iam_execution_role_arn
  service_name            = local.service_name
  subnet_ids              = module.network.default_subnet_ids
  env_task = [
    { name = "ZINC_DATA_PATH", value = var.ZINC_DATA_PATH },
    { name = "ZINC_FIRST_ADMIN_USER", value = var.ZINC_FIRST_ADMIN_USER },
    { name = "ZINC_FIRST_ADMIN_PASSWORD", value = var.ZINC_FIRST_ADMIN_PASSWORD },
    { name = "ZINC_MAX_DOCUMENT_SIZE", value = var.ZINC_MAX_DOCUMENT_SIZE }
  ]
}
#####################

### resourcesc S3 ###
resource "aws_s3_bucket" "data_bucket_name" {
  bucket = local.data_bucket_name
}

resource "aws_s3_bucket_acl" "this" {
  bucket = aws_s3_bucket.data_bucket_name.id
  acl    = "private"
}
#################

### resources Lambda ###
resource "aws_lambda_function" "index_lambda_func" {
  depends_on = [
    module.iam
  ]

  function_name = local.lambda_name
  filename      = data.archive_file.dummy_lambda.output_path
  role          = module.iam["${local.lambda_role_name}"].iam_execution_role_arn
  handler       = local.lambda_handler
  runtime       = local.lambda_runtime
  environment {
    variables = {
      "IDX_ENRONTGZ_FIELDS"     = var.IDX_ENRONTGZ_FIELDS
      "IDX_ENRONTGZ_FORMAT"     = var.IDX_ENRONTGZ_FORMAT
      "IDX_ENRONTGZ_IDXNAME"    = var.IDX_ENRONTGZ_IDXNAME
      "IDX_ENRONTGZ_SEPARATOR"  = var.IDX_ENRONTGZ_SEPARATOR
      "IDX_ENRONTGZ_TERMINATOR" = var.IDX_ENRONTGZ_TERMINATOR
      "IDX_ENRONTGZ_URLZINC"    = "http://${module.network.lb_dns_name}/api/"
      "IDX_ENRONTGZ_TOKENZINC"  = base64encode("${var.ZINC_FIRST_ADMIN_USER}:${var.ZINC_FIRST_ADMIN_PASSWORD}")
    }
  }
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.index_lambda_func.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.data_bucket_name.arn
}

resource "aws_s3_bucket_notification" "index_lambda_notify" {
  depends_on = [
    aws_lambda_permission.allow_bucket
  ]

  bucket = aws_s3_bucket.data_bucket_name.id
  lambda_function {
    lambda_function_arn = aws_lambda_function.index_lambda_func.arn
    events              = ["s3:ObjectCreated:Put"]
  }
}
#################
