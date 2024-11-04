USE testdb;

CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100),
    department VARCHAR(100)
);

INSERT INTO employees (name, email, department) VALUES
    ('John Doe', 'john@example.com', 'Engineering'),
    ('Jane Smith', 'jane@example.com', 'Marketing'),
    ('Bob Wilson', 'bob@example.com', 'Sales'),
    ('Alice Brown', 'alice@example.com', 'Engineering');

CREATE TABLE departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    budget DECIMAL(15,2)
);

INSERT INTO departments (name, budget) VALUES
    ('Engineering', 1000000.00),
    ('Marketing', 500000.00),
    ('Sales', 750000.00);