CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    name TEXT,
    author TEXT,
    genre TEXT,
    description TEXT,
    file_url TEXT
);

INSERT INTO books (name, author, genre, description, file_url) VALUES
('1984', 'George Orwell', 'Dystopia', 'Dystopian novel about surveillance state', '/files/book1.pdf'),
('Animal Farm', 'George Orwell', 'Dystopia', 'Political satire about farm animals', '/files/book2.pdf'),
('War and Peace', 'Leo Tolstoy', 'Fiction', 'Epic historical novel', '/files/book3.pdf'),
('Crime and Punishment', 'Fyodor Dostoevsky', 'Fiction', 'Psychological drama', '/files/book4.pdf'),
('The Brothers Karamazov', 'Fyodor Dostoevsky', 'Philosophy', 'Philosophical novel', '/files/book5.pdf'),
('Clean Code', 'Robert C. Martin', 'Programming', 'Software engineering principles', '/files/book6.pdf'),
('The Pragmatic Programmer', 'Andrew Hunt', 'Programming', 'Programming best practices', '/files/book7.pdf'),
('Go in Action', 'William Kennedy', 'Programming', 'Go language guide', '/files/book8.pdf'),
('Design Patterns', 'Erich Gamma', 'Programming', 'Classic design patterns', '/files/book9.pdf'),
('Harry Potter', 'J.K. Rowling', 'Fantasy', 'Wizarding world story', '/files/book10.pdf');