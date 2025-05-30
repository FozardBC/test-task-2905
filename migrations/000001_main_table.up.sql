CREATE TABLE quotes(
    id SERIAL,
    quote TEXT NOT NULL,
    author TEXT NOT NULL
);

CREATE UNIQUE INDEX people_id_uindex ON quotes (id);