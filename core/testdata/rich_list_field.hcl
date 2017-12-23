resource "simple_resource" "test" {
    rich_list_field {
        nested_field = "One"
        nested_field = "Two"
    }

    rich_list_field {
        nested_field = "Three"
        nested_field = "Four"
    }
}