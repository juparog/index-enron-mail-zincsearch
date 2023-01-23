# main.tf

### Data ###
data "aws_iam_policy_document" "this" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = var.identifiers
    }
  }
}
#############

### Resources ###
resource "aws_iam_role" "this" {
  name               = var.role_name
  assume_role_policy = data.aws_iam_policy_document.this.json
}

resource "aws_iam_role_policy_attachment" "ecsTaskExecutionRole_policy" {
  for_each = toset(var.arn_policies)

  role       = aws_iam_role.this.name
  policy_arn = each.value
}
#################
