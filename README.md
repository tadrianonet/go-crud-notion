 # go-crud-notion

 This is an example project built with Clean Architecture in Go.

 ## How to Run

 1. **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/go-crud-notion.git
    cd go-crud-notion
    ```

 2. **Install dependencies**:
    ```bash
    go mod tidy
    ```

 3. **Run the project**:
    ```bash
    go run cmd/main.go
    ```


 ## Project Structure

 The project follows Clean Architecture, with the following layers:

 - **Entities**: Contains the core business rules.
 - **Use Cases**: Implements specific business rules.
 - **Interfaces**: External layer that handles input/output (HTTP handlers, database).
 - **Repositories**: Implements data persistence.

 ## Pre-Commit Hook

 This project includes a pre-commit hook to ensure code quality before committing changes. It performs the following checks:

 - **Code formatting**: Ensures the code is properly formatted using `gofmt`.
 - **Linting**: Runs `golangci-lint` to catch common issues.
 - **Tests**: Executes all tests in the project using `go test`.

 To use the pre-commit hook, make sure `golangci-lint` is installed:
 ```bash
 go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
 ```

 ## How to Contribute

 1. Fork the project.
 2. Create a branch for your feature (`git checkout -b feature/new-feature`).
 3. Commit your changes (`git commit -m 'Add new feature'`).
 4. Push to the branch (`git push origin feature/new-feature`).
 5. Open a Pull Request.

 ## License

 This project is licensed under the [MIT License](LICENSE).