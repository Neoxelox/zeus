# zeus

**`Unire.one main backend`**

## Setup

- Install:
  - [`golang 1.16 >=`](https://golang.org/dl/)
  - [`docker 20.10.5 >=`](https://docs.docker.com/get-docker/)
  - [`docker-compose 1.28.0 >=`](https://docs.docker.com/compose/install/)
  - [`Python 3.8.6 >=`](https://www.python.org/downloads/)
  - [`PyInvoke 1.5.0 >=`](https://www.pyinvoke.org/installing.html)

Type `invoke --list` for further commands and `invoke --help <command>` for extra information of the command.

**Head to [documentation](https://google.com) for a full explanation of the project.**

## Architecture

The project uses a subset of the Clean Architecture, composed of 3 different layers: **Payload**->**Entity**->**Model**, which are disjoint and the dependencies are unidirectional towards the database domain.

Use Case handlers are created at launch time, meaning that only a single DB connection pool is created. It will be used to create all kind of repositories, that will then be injected to each **Use Case**. That means that **Repositories** and **Handlers** must be thread safe to attend different requests.

Regarding to tests, you should emphasize on unit tests in the **Entity** domain, and integration tests in the **Payload** layer. Repository tests are welcomed, but are less "compulsory". Mocks must be created for every entity methods, handler or repository, so that your tests don't rely on imported packages.

## Structure

Follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

### `/assets`

Static files to be served by the application.

### `/chart`

Configuration files and templates for the Helm chart.

### `/cmd`

Main entrypoint for the microservice. Just an small `main` function that invokes code from `/internal` or `/pkg` and starts the application.

### `/docs`

Design and user documents (in addition to godoc generated documentation).

### `/internal`

Private application code that other projects won't import.

### `/pkg`

Public library code that other external projects could import. All packages inside this folder should be able to be imported in a go program through `go get yourRepository/pkg/yourPkg`. Each package should have:

- `README.md`
- `LICENSE`
- `go.mod` ( create a go module: `go mod init yourPkg` )
- `yourPkg.go` ( use an `init()` function to initialize the package )
- `yourPkg_test.go`
- `yourPkg_mock.go` ( the [testify](github.com/stretchr/testify) package is a good option for making mocks )
- `yourPkg_doc.go` ( optional, for [GoDoc](https://godoc.org) completition )

Notice that your public package **cannot** import code from your private code, thus Go won't be able to compile code imported from an _external_ `/internal` directory.

### `/scripts`

Scripts to perform various build, install, analysis... operations. These scripts keep the root level Makefile/Pyinvoke small and simple.

## License

This project is licensed under the [GPL-3.0 License](https://opensource.org/licenses/GPL-3.0) - read the [LICENSE](LICENSE) file for details.
