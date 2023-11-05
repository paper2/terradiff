resource "null_resource" "expects_create" {
  provisioner "local-exec" {
    command = "echo test"
  }
}
resource "null_resource" "expects_replace" {
  provisioner "local-exec" {
    command = "echo test"
  }
  triggers = {
    trigger = null_resource.expects_create.id
  }
}
resource "null_resource" "changes" {
  provisioner "local-exec" {
    command = "echo test"
  }
  triggers = {
    trigger = null_resource.expects_replace.id
  }
}
