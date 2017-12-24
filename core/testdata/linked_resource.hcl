resource "simple_resource" "test" {
    scalar_field = "${aws_s3_bucket.name.id}"
}