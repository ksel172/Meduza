# Contributing to Meduza

We follow a structured Git branching model:

- **`main`**: The production-ready branch.
- **`dev`**: The development branch.
- **Feature Branches (`feature/<name>`)**: Used for adding new features. Branch from `dev`.
- **Fix Branches (`fix/<name>`)**: For non-urgent fixes (e.g., UI updates, code styling, optimizations). Branch from `dev`.
- **Hotfix Branches (`hotfix/<name>`)**: For urgent production fixes. Branch from `main`.

## How to Contribute

1. **Fork the Repository**: Make a fork of the repository.
2. **Clone Your Fork**: Clone the repository to your machine.
    ```bash
    git clone https://github.com/<your-username>/Meduza.git
    ```
3. **Create a Feature/Fix/Hotfix Branch**: Create a new branch for your changes.
    ```bash
    # For features
    git checkout -b feature/<name> dev

    # For fixes
    git checkout -b fix/<name> dev

    # For hotfixes
    git checkout -b hotfix/<name> main
    ```
4. **Make Changes**: Implement your changes to the new branch.
5. **Commit Your Changes**: Ensure your commit message is clear and concise.
    ```bash
    # For features
    git commit -m "Add feature X"

    # For fixes
    git commit -m "Fix for X"

    # For hotfixes
    git commit -m "Hotfix for X"
    ```
6. **Push Your Changes**: Push the feature branch to your fork.
    ```bash
    git push origin <feature/fix/hostfix>/<name>
    ```
7. **Create a Pull Request**: Open a pull request from your feature branch to `dev`.

## Code Style

Follow the code style of the repository (e.g., Go idioms, file structure etc.).

## Reporting Issues

If you encounter any bugs or have suggestions for improving the project, please check out the existing [issues](https://github.com/ksel172/Meduza/issues) or open a new issue with a clear name and description.

## License

By contributing, you agree that your contributions will be licensed under the project's [BSD-3-Clause License](LICENSE).
