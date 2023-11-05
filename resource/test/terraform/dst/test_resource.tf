resource "null_resource" "expects_destroy" {
  provisioner "local-exec" {
    command = "echo test"
  }
}
resource "null_resource" "expects_replace" {
  provisioner "local-exec" {
    command = "echo test"
  }
}
// NOTE: Plan時にchangesが残るようにする
resource "null_resource" "changes" {
  provisioner "local-exec" {
    command = "echo test"
  }
  triggers = {
    trigger = null_resource.expects_replace.id
  }
}
