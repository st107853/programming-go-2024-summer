CREATE TABLE IF NOT EXISTS users (
    guids VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    ip VARCHAR(128) NOT NULL,
    PRIMARY KEY (guids)
);

INSERT INTO users
    (guids, email, ip)
VALUES
    ('c15dcfee-adfe-4874-9e59-73d54335d214', 'johndoe@gm.com', '217.197.2.10'),
    ('8415c0f5-bf40-4ab2-a230-64db5b080f45', 'johndoe@gm.com', '217.197.2.10'),
    ('6cc105f6-846e-408e-9a79-35008aeb2ddd', 'johndoe@gm.com', '217.197.2.10'),
    ('c5d78c18-d2ae-4608-b456-18bf0aecc04d', 'johndoe@gm.com', '217.197.2.10');