# Note Service

A simple note service with a backend in Golang and a frontend in native HTML, CSS, and JavaScript.

## Project Structure
```bash
note-service/
├── frontend/
│ ├── css/
│ │ └── styles.css
│ ├── js/
│ │ └── scripts.js
│ └── index.html
├── handlers/
│ └── note.go
├── models/
│ └── note.go
├── routers/
│ └── router.go
├── storage/
│ ├── memory.go
│ └── storage.go.todo
├── utils/
│ └── response.go
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```
## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (1.16+)
- A web browser

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/f0rxz/note-service.git
    cd note-service
    ```

2. Install the dependencies:

    ```bash
    go mod tidy
    ```

### Running the Application

1. Start the backend server:

    ```bash
    go run main.go
    ```

    The server will be running at `http://127.0.0.1`.

2. Open your web browser and navigate to `http://127.0.0.1` to access the frontend.

### API Endpoints

- `GET /api/notes`: Retrieve all notes
- `POST /api/notes`: Create a new note
- `GET /api/notes/{id}`: Retrieve a specific note by ID
- `PUT /api/notes/{id}`: Update a specific note by ID
- `DELETE /api/notes/{id}`: Delete a specific note by ID

### Frontend

The frontend is a simple HTML, CSS, and JavaScript application. It provides an interface to:

- View all notes
- Create a new note
- Delete a note

### Project Implementation

#### `main.go`

Sets up and starts the server.

#### `routers/router.go`

Configures the router, defines API routes, and serves static files.

#### `handlers/note.go`

Contains handlers for API endpoints.

#### `models/note.go`

Defines the `Note` model.

#### `storage/memory.go`

Implements in-memory storage for notes.

#### `storage/storage.go.todo`

A placeholder for additional storage backends. This file can be used to implement different storage mechanisms, such as database storage or file-based storage, to replace or augment the in-memory storage. It provides a structure for future expansion of the storage capabilities of the application.

#### `utils/response.go`

Utility functions for responding with JSON.

### Frontend Files

- `frontend/index.html`: The main HTML file.
- `frontend/css/styles.css`: The CSS file for styling.
- `frontend/js/scripts.js`: The JavaScript file for client-side logic.

### Contributing

Feel free to open issues or submit pull requests for any improvements.

### License

This project is licensed under the MIT License. See the `LICENSE` file for details.
