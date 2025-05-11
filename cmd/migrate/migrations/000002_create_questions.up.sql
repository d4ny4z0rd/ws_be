CREATE TABLE IF NOT EXISTS dsa_questions (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    input_format TEXT NOT NULL,
    output_format TEXT NOT NULL,
    example_input TEXT NOT NULL,
    example_output TEXT NOT NULL
);
