<p align="center">
<img src="https://github.com/andygeiss/cloud-native-store/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Store

[![License](https://img.shields.io/github/license/andygeiss/cloud-native-store)](https://github.com/andygeiss/cloud-native-store/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-store)](https://github.com/andygeiss/cloud-native-store/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/cloud-native-store)](https://goreportcard.com/report/github.com/andygeiss/cloud-native-store)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/df82f7d9fa07469eadd726342e837197)](https://app.codacy.com/gh/andygeiss/cloud-native-store/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/df82f7d9fa07469eadd726342e837197)](https://app.codacy.com/gh/andygeiss/cloud-native-store/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

**Cloud Native Store** is a Go-based key-value store showcasing cloud-native patterns like transactional logging, data sharding, encryption, TLS, and circuit breakers. Built with hexagonal architecture for modularity and extensibility, it includes a robust API and in-memory storage for efficiency and stability.

## Project Motivation

The motivation behind **Cloud Native Store** is to provide a practical example of implementing a key-value store that adheres to cloud-native principles. The project aims to:

1. Demonstrate best practices for building scalable, secure, and reliable cloud-native applications.
2. Showcase the use of hexagonal architecture to enable modular and testable code.
3. Offer a reference implementation for features like encryption, transactional logging, and stability mechanisms.
4. Inspire developers to adopt cloud-native patterns in their projects.

## Project Setup and Run Instructions

Follow these steps to set up and run the **Cloud Native Store**:

### Prerequisites
1. Create an encryption key by running the following command:
```bash
just genkey
```

2. Create an `.env` file and replace the following values besides `HOME_PATH` with your own:

```env
ENCRYPTION_KEY="0a0375de7bd186c2f8d80ef94e5f3d357462f594ca6785d4779f52bcb2b65b85"
GITHUB_CLIENT_ID=""
GITHUB_CLIENT_SECRET=""
GITHUB_REDIRECT_URL="http://localhost:8080/auth/callback"
GITHUB_SCOPE="user:read"
HOME_PATH="/ui/store"
```

### Commands

#### Test the Service
To run unit tests:
```bash
just test
```
This will execute tests for the core service logic.

#### Run the Service
To start the service:
```bash
just run
```

#### How to Test

After running the service, you can verify its health by visiting the UI in your browser:

[http://localhost:8080/ui](http://localhost:8080/ui)
