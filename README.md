# book-manager-service

In this mini-project, we want to implement a small service for maintaining users’ books. With this service, clients
can create new user accounts. The requests of book-related API should be authenticated.

## Structure of packages

* database package which is responsible for the CRUD operations.
* handlers package which contains the HTTP handlers of the project.
* authenticate package which handles the JWT authentication.
* config package which contains the configuration of this backend service.

## APIs of this service

### Signup

path: /api/v1/auth/signup <br>
<strong>Method</strong>: POST

This API creates a new user. Here is the model of the users

```go
// Request
{
    "user_name": "username",
    "email": "email@gmail.com",
    "password": "password",
    "first_name": "fname",
    "last_name": "lname",
    "phone_number": "09136594879",
    "gender": "gender"
}
```

### Login

path: /api/v1/auth/login <br>
<strong>Method</strong>: POST

This API authenticates the user and returns the access token for the user.

```go
// Request
{
    "username": "username",
    "password": "password"
}
```
```go
// Response
{
    "access_token": "abc.abc.abc"
}
```

### CreateBook

path: /api/v1/books <br>
<strong>Method</strong>: POST

This API creates a new book for logged-in users.

```go
// Request Body
{
    "name": "book_name",
    "author": {
        "first_name": "fname of author",
        "last_name": "lname of author",
        "birthday": "2000-01-12T00:00:00+03:30",
        "nationality": "french"
    },
    "category": "book_category",
    "volume": 1,
    "published_at": "2000-01-12T00:00:00+03:30",
    "summary": "this is a summary of the book.",
    "table_of_contents": [
        "fasle_1", "fasle_2"
    ],
    "publisher": "publisher_name"
}
```

### GetAllBooks

path: /api/v1/books <br>
<strong>Method</strong>: GET

This API retrieves all books for all users and returns a list of books.

```go
// Response of the /api/v1/books
{
    "books" [
        {
            "name": "book_name",
            "author": "<fname of author> <lname of author>",
            "category": "book_category",
            "volume": 1,
            "published_at": "2000-01-12T00:00:00+03:30",
            "summary": "this is a summary of the book.",
            "publisher": "publisher_name"
        },
        {
            "name": "book_name",
            "author": "<fname of author> <lname of author>",
            "category": "book_category",
            "volume": 1,
            "published_at": "2000-01-12T00:00:00+03:30",
            "summary": "this is a summary of the book.",
            "publisher": "publisher_name"
        }
    ]
}
```

### GetBook

path: /api/v1/books/\<id\> <br>
<strong>Method</strong>: GET

This API retrieves a specific book by id and returns the full information about that book.

```go
// Response Body of /api/v1/books/13
{
    "name": "book_name",
    "author": {
        "first_name": "fname of author",
        "last_name": "lname of author",
        "birthday": "2000-01-12T00:00:00+03:30",
        "nationality": "french"
    },
    "category": "book_category",
    "volume": 1,
    "published_at": "2000-01-12T00:00:00+03:30",
    "summary": "this is a summary of the book.",
    "table_of_contents": [
        "fasle_1", "fasle_2"
    ],
    "publisher": "publisher_name"
}
```

### UpdateBook

path: /api/v1/books/\<id\> <br>
<strong>Method</strong>: UPDATE

This API updates the information about books.

```go
// Request body of the update API /api/v1/books/12
{
    "name": "book_name_updated",
    "category": "book_category_updated",
}
```

### DeleteBook

This API lets authenticated users delete the book they created.