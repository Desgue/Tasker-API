# Tasker Project Management API

**The Frontend repository of this Application can be found [here](https://github.com/Desgue/tasker)**

**The link for the Amplify Hosted website can be found [here](https://production.d3ozduy4s4mqlc.amplifyapp.com/)**



## Table of Contents
1. [Development Setup](#development-setup)
    - [Database](#database)
    - [Dependencies](#dependencies)
    - [Additional Notes](#additional-notes)
2. [Features](#features)
    - [Implemented Features](#implemented-features)
    - [Planned Features](#planned-features)
3. [System Architecture](#system-architecture)
4. [API Endpoints](#api-endpoints)
    - [Projects API](#projects-api)
    - [Tasks API](#tasks-api)
       


## Development Setup

### Database:

- **PostgreSQL:**
  - Uses the official PostgreSQL Docker image for local database management.

    ```bash
    docker run --name <your_name>_db -e POSTGRES_USER=<your_username> -e POSTGRES_PASSWORD=<your_password> -p 5432:5432 postgres
    ```


### Dependencies:

- **JWT Verification:**
  - [lestrrat-go/jwx](https://github.com/lestrrat-go/jwx) validates and verifies JWT tokens issued by Amazon Cognito.

- **Environment Management:**
  - [joho/godotenv](https://github.com/joho/godotenv) parses and utilizes values from .env files.

- **Database Interaction:**
  - [lib/pq](https://github.com/lib/pq) enables interaction with the PostgreSQL database.

- **CORS Handling:**
  - [rs/cors](https://github.com/rs/cors) configures and handles Cross-Origin Resource Sharing (CORS) requests.


### Additional Notes:

- **Development:**
  - This guide is for development environments. Production configurations may differ.

# Features

### Implemented Features

- **Cognito Auth:**
  - Secure user authentication and authorization handled by Amazon Cognito.

- **Project CRUD:**
  - Perform Create, Read, Update, and Delete operations on projects.

- **Task CRUD:**
  - Perform Create, Read, Update, and Delete operations on tasks associated with projects.

### Planned Features

- **Add Teams and Project Collaboration:**
  - Collaborate with teams and friends on projects for enhanced project management.

- **Project Role Assignment:**
  - Assign roles to users for better project organization and collaboration.
 
    
## System Architecture

![System Architecture Diagram](https://github.com/Desgue/tasker/blob/main/public/tasker-diagram2.drawio.svg)

1. User registration or login initiates the process, with information securely stored in the Amazon Cognito user pool.
2. Each user request includes an authentication token sent to the Golang API.
3. The Golang API validates the token with Amazon Cognito.
4. Upon successful validation, the server responds with the requested data.


## API Endpoints

### Projects API

#### GET /projects

**Description:** Retrieves a list of all projects associated with the authenticated user.

**Returned Data:**
- `id`: Unique identifier of the project (integer)
- `title`: Title of the project (string)
- `description`: Brief description of the project (string)
- `priority`: Priority level of the project (string, one of "Low", "Medium", "High")
- `created_at`: Date and time the project was created (ISO 8601 format)

#### GET /projects/{projectId}

**Description:** Retrieves details of a specific project identified by its unique `projectId`.

**Returned Data:**
- `id`: Unique identifier of the project (integer)
- `title`: Title of the project (string)
- `description`: Brief description of the project (string)
- `priority`: Priority level of the project (string, one of "Low", "Medium", "High")
- `created_at`: Date and time the project was created (ISO 8601 format)

#### POST /projects

**Description:** Creates a new project.

**Required Data:**
- `title`: Title of the project (string)
- `description`: Brief description of the project (string)
- `priority`: Priority level of the project (string, one of "Low", "Medium", "High") (Optional, defaults to "Low")

#### PUT /projects/{projectId}

**Description:** Updates an existing project identified by its unique `projectId`.

**Required Data:**
- `title`: Title of the project (string)
- `description`: Brief description of the project (string)
- `priority`: Priority level of the project (string, one of "Low", "Medium", "High")

#### DELETE /projects/{projectId}

**Description:** Deletes a project identified by its unique `projectId`. This action also removes all associated tasks.

### Tasks API

#### GET /projects/{projectId}/tasks

**Description:** Retrieves a list of all tasks associated with a specific project identified by its unique `projectId`.

**Returned Data:**
- `id`: Unique identifier of the task (integer)
- `title`: Title of the task (string)
- `description`: Brief description of the task (string)
- `status`: Current status of the task (string, one of "Pending", "In Progress", "Done")
- `created_at`: Date and time the task was created (ISO 8601 format)

#### GET /projects/{projectId}/tasks/{taskId}

**Description:** Retrieves details of a specific task identified by its unique `taskId` within a project identified by its `projectId`.

**Returned Data:**
- `id`: Unique identifier of the task (integer)
- `title`: Title of the task (string)
- `description`: Brief description of the task (string)
- `status`: Current status of the task (string, one of "Pending", "In Progress", "Done")
- `created_at`: Date and time the task was created (ISO 8601 format)

#### POST /projects/{projectId}/tasks

**Description:** Creates a new task within a project identified by its unique `projectId`.

**Required Data:**
- `title`: Title of the task (string)
- `description`: Brief description of the task (string)
- `status`: Initial status of the task (string, one of "Pending", "In Progress", "Done") (Optional, defaults to "Pending")

#### PUT /projects/{projectId}/tasks/{taskId}

**Description:** Updates an existing task identified by its unique `taskId` within a project identified by its `projectId`.

**Required Data:**
- `title`: Title of the task (string)
- `description`: Brief description of the task (string)
- `status`: Updated status of the task (string, one of "Pending", "In Progress", "Done")

#### DELETE /projects/{projectId}/tasks/{taskId}

**Description:** Deletes a specific task identified by its unique `taskId` within a project identified by its `projectId`.



