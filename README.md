# hamed-bank

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hamedblue1381/hamed-bank)](https://goreportcard.com/report/github.com/Hamedblue1381/hamed-bank)

`Hamed-bank` is a simple and efficient repository that demonstrates the implementation of gRPC in Golang with a PostgreSQL database. This project is designed with a clean and modular architecture, making it a great starting point for developers to build their own gRPC-based applications using Go and PostgreSQL.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [License](#license)

## Features

- Well-structured and modular codebase.
- gRPC implementation for efficient communication between client and server.
- PostgreSQL integration for a robust and scalable database.
- Docker support for easy deployment and development.
- Comprehensive unit tests for enhanced reliability.

## Requirements

- [Go](https://golang.org/doc/install) 1.16 or later
- [PostgreSQL](https://www.postgresql.org/download/) 12 or later
- [Docker](https://www.docker.com/get-started) 20.10 or later (optional)
- [GNU Make](https://www.gnu.org/software/make/) 3.81 or later

## Installation

1. Clone the repository:

   ````bash
   git clone https://github.com/Hamedblue1381/go-postgres-gRPC
   cd go-postgres-gRPC
   ```

2. Install the required Go dependencies:

   ````bash
   go mod download
   ```

## Quick Start

1. Start the PostgreSQL database using Docker:

   ````bash
   make postgres
   ```
2. Create the postgres database:

   ````bash
   make createdb
   ```
3. Migrate up the database using go migration tool:

   ````bash
   make migrateup
   ```

4. Run the gRPC server:

   ````bash
   make server
   ```
   

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
