resource "simple_resource" "test" {
    rich_list_field {
        nested_field = "One"
    }

    rich_list_field {
        nested_field = "Two"
    }
}