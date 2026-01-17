# Go Interview

Go Interview is a Retrieval-Augmented Generation (RAG) system designed to collect user biographies by interviewing them about specific life areas. The application leverages a Domain-Driven Design (DDD) approach to manage the complexity of interviews, biography generation, and the underlying memory systems required for effective RAG operations.

## Architecture Overview

The project follows a modular architecture organized by domains within the `internal/` directory:

### 1. Biography Domain (`internal/biography`)
Focused on the core subject matter of the application.
- **Life Areas**: Defines the different areas of a user's life to be explored.
- **Criteria**: Manages the specific criteria used to evaluate or structure the biographical data.

### 2. Interview Domain (`internal/interview`)
Handles the interaction and data collection process.
- **Transcripts**: Manages the structured records of interviews.
- **Raw Data**: Handles the raw inputs from the interview process.

### 3. Memory Domain (`internal/memory`)
Provides the RAG foundation for the system.
- **Facts**: atomic units of information extracted from interviews.
- **Embeddings & Vectors**: Supports semantic search and retrieval capabilities essential for the RAG system.

## Project Structure

```text
├── cmd/
│   └── app/            # Main application entry point
├── internal/
│   ├── biography/      # Biography domain logic
│   ├── common/         # Shared utilities and domain objects
│   ├── interview/      # Interview domain logic
│   └── memory/         # RAG and vector memory logic
├── go.mod              # Go module definition
└── README.md           # Project documentation
```

## Getting Started

### Prerequisites

- **Go**: Version 1.25.5 or higher.

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-interview
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

### Running the Application

To run the main application entry point:

```bash
go run cmd/app/main.go
```

*Note: Currently, the main application serves as a placeholder and prints a welcome message.*

### Testing

Run the test suite across all packages (once tests are implemented):

```bash
go test ./...
```
