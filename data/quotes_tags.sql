CREATE TABLE IF NOT EXISTS quotes_tags (
    quote_id integer REFERENCES quotes(id),
    tag_id integer REFERENCES tags(id),
    PRIMARY KEY (quote_id, tag_id)
);