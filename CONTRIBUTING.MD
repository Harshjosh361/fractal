# Contributing to Fractal

Thank you for considering contributing to **Fractal**! Contributions are what make the open-source community a wonderful place to learn, inspire, and create. Whether it's fixing a bug, implementing a new feature, improving documentation, or suggesting an enhancement, your help is highly appreciated.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [How to Contribute](#how-to-contribute)
    - [Bug Reports](#bug-reports)
    - [Feature Requests](#feature-requests)
    - [Code Contributions](#code-contributions)
4. [Development Guidelines](#development-guidelines)
5. [Pull Request Process](#pull-request-process)
6. [Contact](#contact)

---

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/0/code_of_conduct.html). By participating, you are expected to uphold this code. Please report unacceptable behavior to `maintainers@fractal.dev`.

---

## Getting Started

To start contributing:
1. Fork the repository.
2. Clone your forked repository locally:
   ```bash
   git clone https://github.com/SkySingh04/fractal.git
   ```
3. Set up the project by installing dependencies:
   ```bash
   cd fractal
   go mod tidy
   ```
4. Create a new branch for your work:
   ```bash
   git checkout -b feature/<feature-name>
   ```

---

## How to Contribute

### Bug Reports
If you find a bug:
- **Check existing issues** to ensure it hasn't already been reported.
- Submit a detailed issue, including:
  - Steps to reproduce the bug.
  - Expected vs. actual behavior.
  - Environment details (OS, Go version, etc.).

### Feature Requests
If you have an idea for a new feature:
- **Open a feature request** on the [Issues page](https://github.com/SkySingh04/fractal/issues).
- Clearly describe:
  - The problem the feature solves.
  - The proposed solution.

### Code Contributions
1. Find an issue labeled `help wanted` or `good first issue` to start.
2. Comment on the issue to let maintainers know you are working on it.
3. Follow the development guidelines below to implement the changes.

---

## Development Guidelines

### Code Style
- Follow the Go community’s [style guide](https://golang.org/doc/effective_go).
- Use `go fmt` to format your code.

### Project Structure
The project uses a modular structure:
- **`cmd/`**: Main entry points for CLI and server modes.
- **`integrations/`**: Data source and destination integrations.
- **`internal/`**: Core functionalities like transformations, validations, etc.
- **`pkg/`**: Shared utilities and libraries.
- **`config/`**: YAML-based configuration handling.

### Running Tests
Run the test suite before submitting a pull request:
```bash
go test ./...
```

If adding new functionality, include unit tests to verify your changes.

### Adding Integrations
To add a new integration (input/output):
1. Follow the [Integration Guide](#integration-guide) in the main documentation.
2. Add test cases in the `integrations` directory.

---

## Pull Request Process

1. **Sync your fork**: Ensure your fork is up-to-date with the main repository:
   ```bash
   git fetch upstream
   git merge upstream/main
   ```
2. **Commit changes**: Follow commit message conventions:
   - Use clear and descriptive messages (e.g., `fix: resolve data validation issue`).
   - Use multiple commits for distinct changes.
   ```bash
   git add .
   git commit -m "feat: add RabbitMQ integration"
   ```
3. **Push changes**:
   ```bash
   git push origin feature/<feature-name>
   ```
4. **Create a pull request**:
   - Go to the main repository.
   - Navigate to the **Pull Requests** tab and click "New Pull Request."
   - Select your branch and explain your changes in detail.

5. **Respond to Feedback**:
   - Be prepared to revise your code based on maintainers' feedback.

---

## Contact

For any questions, feel free to reach out:
- Open an issue for general queries or bug reports.

Thank you for contributing to **Fractal**! Together, we can build a powerful and flexible data processing tool. 🎉
--- 

