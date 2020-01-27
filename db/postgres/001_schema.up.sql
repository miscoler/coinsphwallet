CREATE TABLE accounts (id text PRIMARY KEY NOT NULL, balance numeric NOT NULL, currency TEXT NOT NULL);
CREATE TABLE transactions (id BIGSERIAL PRIMARY KEY NOT NULL, sender TEXT, receiver TEXT, amount numeric NOT NULL, currency TEXT NOT NULL);

INSERT INTO accounts (id, balance, currency) VALUES ('johny guitar', '67.0', 'USD');
INSERT INTO accounts (id, balance, currency) VALUES ('courier', '100.0', 'CAPS');
INSERT INTO accounts (id, balance, currency) VALUES ('mr house', '10000.0', 'CAPS');
INSERT INTO accounts (id, balance, currency) VALUES ('yes man', '10.0', 'CAPS');
INSERT INTO accounts (id, balance, currency) VALUES ('caesar', '1000.0', 'DENARIUS');

