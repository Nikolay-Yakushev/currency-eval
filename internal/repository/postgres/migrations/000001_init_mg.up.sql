
CREATE TABLE currencies (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            value FLOAT NOT NULL,
                            date TIMESTAMP NOT NULL
);
