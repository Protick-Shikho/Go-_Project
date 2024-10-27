# To-do Task Management Application

This is a simple task management application built with Go and MySQL. Users can create, read, update, and delete tasks in a to-do list format.

## Features

- Create a new task
- Show all tasks
- Update task status
- Delete a task

## Prerequisites

- Go 1.16 or later
- MySQL database

## Installation

1. **Clone the repository:**

   ```bash
   git clone git@github.com:Protick-Shikho/To_do_Task.git
   cd To_do_Task

2. Set up the MySQL database:

    Open your MySQL shell or a database management tool.

    Run the following commands to create the database and table:

CREATE DATABASE todo_list;

USE todo_list;

CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status ENUM('pending', 'completed') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

3. Create a .env file:

In the root directory of the project, create a .env file with your MySQL credentials:


DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=todo_list
DB_HOST=127.0.0.1
DB_PORT=3306

Replace your_password_here with your actual MySQL root password.


4. Run The File
