# Git Implementation in Go

[![progress-banner](https://backend.codecrafters.io/progress/git/211abd6c-a869-410a-a9e1-512a84007df8)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This project is a Go implementation of core Git functionality, built as part of the [Codecrafters "Build your own Git" challenge](https://app.codecrafters.io/courses/git/overview). This implementation provides a hands-on understanding of Git's internal workings and data structures.

## 📌 What This Project Does

This Git implementation covers the fundamental operations that Git performs under the hood:

- **Repository Initialization:** Creating a new Git repository with proper directory structure
- **Object Storage:** Storing and retrieving Git objects (blobs, trees, commits) using SHA-1 hashing
- **Tree Operations:** Creating and reading Git tree objects that represent directory structures
- **Commit Management:** Creating commit objects with proper parent relationships
- **Content Addressing:** Using SHA-1 hashes to uniquely identify and retrieve objects
- **Compression:** Storing objects efficiently using zlib compression

_Note: I haven't implemented remote cloning yet, but it's coming next._

## ✨ Key Features

- **Git Object Model:** Full implementation of Git's object storage system with support for:
  - **Blobs:** Store file contents
  - **Trees:** Store directory structures and file metadata
  - **Commits:** Store commit information with parent relationships
- **SHA-1 Hashing:** Content-addressable storage using SHA-1 hashes for object identification
- **Compression:** Objects are compressed using zlib before storage
- **Working Directory Operations:** Create trees from current directory structure
- **Object Inspection:** Display and analyze Git objects
- **Reference Management:** Basic handling of Git references and object relationships

## 🛠️ Why I Built This Project

As a developer who uses Git daily, I wanted to understand what happens behind the scenes when I run commands like `git add`, `git commit`, and `git init`. This project served as an excellent opportunity to:

- **Demystify Git's Magic:** Understand how Git stores data and manages versions
- **Learn Systems Programming:** Work with binary formats, compression, and file systems
- **Practice Go Development:** Gain deeper experience with Go's standard library and error handling
- **Explore Data Structures:** Understand how Git's DAG (Directed Acyclic Graph) structure works
- **Build Foundation Knowledge:** Prepare for implementing more complex Git operations like cloning

## 🔍 How It Works Internally

The Git implementation processes operations through several key components:

#### 1. **Object Storage System**

- All Git objects are stored in objects directory
- Objects are compressed using zlib and stored in subdirectories based on their SHA-1 hash
- Each object has a header format: `<type> <size>\0<content>`

#### 2. **Command Processing**

- The `GitClient` routes commands to appropriate executors
- Each command (init, hash-object, cat-file, etc.) has its own implementation
- Commands follow a consistent interface pattern for easy extension

#### 3. **Tree and Commit Operations**

- Trees store directory structures with file modes, names, and hashes
- Commits reference a tree hash and contain metadata (author, message, parent)
- The implementation properly handles the binary format used by Git

#### 4. **Object Hashing and Storage**

- Content is hashed using SHA-1 algorithm
- Objects are compressed with zlib before storage
- Directory structure follows Git's standard: first two characters of hash as directory name

## ⚙️ How to Set Up and Run

### Prerequisites

- Go 1.24 or later (as specified in `go.mod`)

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/md-talim/codecrafters-git-go.git
   cd codecrafters-git-go
   ```

2. Build the project:
   ```sh
   go build -o mygit app/main.go
   ```

### Running the Git Implementation

**General Usage:**

```sh
./mygit <command> [<args>...]
```

**Available Commands:**

| Command       | Description                                 | Example                                                         |
| ------------- | ------------------------------------------- | --------------------------------------------------------------- |
| `init`        | Initialize a new Git repository             | `./mygit init`                                                  |
| `hash-object` | Compute object hash, optionally store it    | `./mygit hash-object [-w] <file>`                               |
| `cat-file`    | Display contents of a Git object            | `./mygit cat-file -p <hash>`                                    |
| `ls-tree`     | List contents of a tree object              | `./mygit ls-tree [--name-only] <tree-hash>`                     |
| `write-tree`  | Create a tree object from current directory | `./mygit write-tree`                                            |
| `commit-tree` | Create a commit object                      | `./mygit commit-tree <tree-hash> -p <parent-hash> -m <message>` |

_Note: `clone` command is planned for future implementation._

### Example Usage

1. **Initialize a repository:**

   ```sh
   ./mygit init
   ```

2. **Create and store a blob:**

   ```sh
   echo "Hello, Git!" > hello.txt
   ./mygit hash-object -w hello.txt
   ```

3. **Examine the stored object:**

   ```sh
   ./mygit cat-file -p <hash-from-previous-command>
   ```

4. **Create a tree from current directory:**

   ```sh
   ./mygit write-tree
   ```

5. **List tree contents:**

   ```sh
   ./mygit ls-tree <tree-hash>
   ```

6. **Create a commit:**
   ```sh
   ./mygit commit-tree <tree-hash> -m "Initial commit"
   ```

## 💡 What I Learned

### Git Internals

- **Object Model:** Git stores everything as objects (blobs, trees, commits) identified by SHA-1 hashes
- **Content Addressing:** Files with identical content have the same hash, enabling efficient storage
- **Tree Structure:** Git represents directories as tree objects containing references to blobs and other trees
- **Commit Graph:** Commits form a directed acyclic graph through parent relationships
- **Object Storage:** How Git organizes objects in the filesystem using hash-based directory structure
- **Binary Formats:** Understanding Git's internal binary formats for storing tree and commit data

### Systems Programming in Go

- **Binary Data Handling:** Working with byte slices, binary encoding, and bit manipulation
- **Error Handling:** Proper error propagation and handling in complex operations
- **File System Operations:** Creating directory structures and managing file permissions
- **Compression:** Using zlib for efficient object storage
- **Hash Functions:** Implementing SHA-1 hashing for content addressing

## 🎯 Challenges & Solutions

- **Binary Format Handling:** Git's tree format uses null-terminated strings and binary hashes requiring careful parsing
- **Compression Management:** Implementing proper zlib compression/decompression while maintaining data integrity
- **Hash Calculations:** Ensuring SHA-1 hashing matches Git's exact format including headers
- **File System Abstractions:** Creating clean interfaces for object storage that mirror Git's behavior
- **Cross-Platform Compatibility:** Ensuring file permissions and paths work correctly across different operating systems

## 🚀 Future Enhancements

The next major milestone is implementing the `clone` command, which will involve:

- **Smart HTTP Protocol:** Implementing Git's transfer protocol for remote repositories
- **Pack File Support:** Handling Git's packfile format for efficient object transfer
- **Network Operations:** HTTP client implementation for repository discovery and data transfer
- **Working Directory Checkout:** Recreating files and directories from downloaded objects
- **Reference Management:** Proper handling of remote refs and branch setup

## 🙏 Acknowledgments

- This project implements the Git version control system as described in the [Pro Git book](https://git-scm.com/book)
- The [Codecrafters "Build your own Git" challenge](https://app.codecrafters.io/courses/git/overview) provided the structure and test cases
- Git's documentation and source code served as invaluable references for understanding the internal
