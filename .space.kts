job("Test") {
    container(displayName = "Test", image = "golang") {
        shellScript {
            content = """
                go version
                go vet .
                go test ./test
            """
        }
    }
}