# Go-Todo-REST-API

A RESTful API example for simple todo application with Go. Idea from [mingrammer/go-todo-rest-api-example](https://github.com/mingrammer/go-todo-rest-api-example). 

Uses [GORM](https://gorm.io/) as ORM library and [gorilla/mux](https://github.com/gorilla/mux) for routing.

## API

#### /projects

* `GET` : Get all projects
* `POST` : Create a new project
  
  #### /projects/:title
* `GET` : Get a project
* `PUT` : Update a project
* `DELETE` : Delete a project
  
  #### /projects/:title/archive
* `PUT` : Archive a project
* `DELETE` : Restore a project 
  
  #### /projects/:title/tasks
* `GET` : Get all tasks of a project
* `POST` : Create a new task in a project
  
  #### /projects/:title/tasks/:id
* `GET` : Get a task of a project
* `PUT` : Update a task of a project
* `DELETE` : Delete a task of a project
  
  #### /projects/:title/tasks/:id/complete
* `PUT` : Complete a task of a project
* `DELETE` : Undo a task of a project



## Todo

- [x] Basic REST API

- [x] Use gorm and gorilla/mux

- [x] (Integration) Tests

- [ ] Authentification

- [ ] Dockerize

- [ ] GoDoc documentation

- [ ] ...