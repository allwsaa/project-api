# Project-Management api

A service for managing projects, tasks, and users. Supports CRUD operations for users, tasks, and projects.


## RENDER URL: https://project-api-xentvlbl.onrender.com/swagger/index.html


## API Endpoints

### Users

- **GET /users**: Get all users.
- **POST /users**: Create a new user.
- **GET /users/{id}**: Get a user by ID.
- **PUT /users/{id}**: Update a user by ID.
- **DELETE /users/{id}**: Delete a user by ID
- **GET /users/{id}/tasks**: Get tasks assigned to a user.
- **GET /users/search?name={name}**: Search users by name.
- **GET /users/search?email={email}**: Search users by email.
  
### Tasks


-  **GET /tasks**: Get all tasks.
- **POST /tasks**: Create a new task.
- **GET /tasks/{id}**: Get a task by ID.
- **PUT /tasks/{id}**: Update a task by ID.
- **DELETE /tasks/{id}**: Delete a task by ID.
- **GET /tasks/search?title={title}**: Search tasks by title.
- **GET /tasks/search?status={status}**: Search tasks by status.
- **GET /tasks/search?priority={priority}**: Search tasks by priority.
- **GET /tasks/search?assignee={userId}**: Search tasks by assignee.
- **GET /tasks/search?project={projectId}**: Search tasks by project.
 
### Projects

- **GET /projects**: Get all projects.
- **POST /projects**: Create a new project.
- **GET /projects/{id}**: Get a project by ID.
- **PUT /projects/{id}**: Update a project by ID.
- **DELETE /projects/{id}**: Delete a project by ID.
- **GET /projects/{id}/tasks**: Get tasks in a project.
- **GET /projects/search?title={title}**: Search projects by title.
- **GET /projects/search?manager={userId}**: Search projects by manager.
  

## HTTP Responses

- **200**: Successful GET, PUT, DELETE requests.
- **201**: Successful POST requests.
- **400**: Invalid request.
- **404**: Resource not found.
- **405**: Method not allowed.


