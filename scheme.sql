CREATE TABLE comments(
                         id SERIAL PRIMARY KEY,
                         newsid BIGINT,
                         comment TEXT,
                         dateunix BIGINT,
                         parents BIGINT,
                         allow bool
)