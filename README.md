# zeus

**`Unire.one main backend`**

## Setup

- Install:
  - [`golang 1.16 >=`](https://golang.org/dl/)
  - [`docker 20.10.5 >=`](https://docs.docker.com/get-docker/)
  - [`docker-compose 1.28.0 >=`](https://docs.docker.com/compose/install/)
  - [`Python 3.8.6 >=`](https://www.python.org/downloads/)
  - [`PyInvoke 1.5.0 >=`](https://www.pyinvoke.org/installing.html)
  - [`Dotenv 0.17.0 >=`](https://pypi.org/project/python-dotenv/)

Type `invoke --list` for further commands and `invoke --help <command>` for extra information of the command.

**Head to [documentation](https://google.com) for a full explanation of the project.**

## Architecture

The project uses a subset of the Clean Architecture, composed of 3 different layers: **Handler**->**Use Case**->**Repository**, which are disjoint and the dependencies are unidirectional towards the repository domain.

Use Case handlers are created at launch time, meaning that only a single DB connection pool is created. It will be used to create all kind of repositories, that will then be injected to each **Use Case**. That means that **Repositories** and **Handlers** must be thread safe to attend different requests.

Regarding to tests, you should emphasize on unit tests in the **Use Case** domain, and integration tests in the **Handler** layer. **Repository** domain tests are welcomed, but are less "compulsory". Mocks must be created for every use case or repository, so that your tests don't rely on imported packages.

## Structure

Follows a subset of the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

### `/assets`

Static files to be served by the application.

### `/chart`

Configuration files and templates for the Helm chart.

### `/cmd`

Main entrypoint for the microservice. Just an small `main` function that invokes code from `/internal` and starts the application.

### `/docs`

Design and user documents (in addition to godoc generated documentation).

### `/internal`

Private application code that other projects won't import.

### `/migrations`

Configuration files and templates for the migrations.

### `/pkg`

Architecture logic of the microservice, that is handlers, payloads, models, use cases and repositories.

### `/scripts`

Scripts to perform various build, install, analysis... operations. These scripts keep the root level Makefile/Pyinvoke small and simple.

## License

This project is licensed under the [GPL-3.0 License](https://opensource.org/licenses/GPL-3.0) - read the [LICENSE](LICENSE) file for details.
