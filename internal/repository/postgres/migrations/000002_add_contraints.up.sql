ALTER TABLE currencies ADD CONSTRAINT unique_name_date UNIQUE (name, date);
