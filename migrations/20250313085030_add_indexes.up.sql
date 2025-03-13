CREATE INDEX idx_books_updated_at_desc ON books (updated_at DESC);
CREATE INDEX idx_authors_updated_at_desc ON authors (updated_at DESC);
CREATE INDEX idx_book_authors_book_id ON book_authors (book_id);