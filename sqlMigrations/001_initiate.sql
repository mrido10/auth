-- +migrate Up
-- +migrate StatementBegin

CREATE SEQUENCE IF NOT EXISTS access_pkey_seq;
CREATE SEQUENCE IF NOT EXISTS user_account_pkey_seq;

CREATE TABLE IF NOT EXISTS access (
  id            BIGINT NOT NULL DEFAULT nextval('access_pkey_seq'::regclass),
  access_name   VARCHAR(50),
  access_code   VARCHAR(50),
  created       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT    access_pkey_seq_id PRIMARY KEY (id)
);

INSERT INTO access (access_name, access_code)
VALUES ('Money Manager', 'mon-mgr');

CREATE TABLE IF NOT EXISTS user_account (
    id          BIGINT NOT NULL DEFAULT nextval('user_account_pkey_seq'::regclass),
    email       VARCHAR(50) UNIQUE,
    password    VARCHAR(100),
    name        VARCHAR(50),
    gender      VARCHAR(10),
    acces_id    BIGINT NOT NULL,
    created     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    modify      TIMESTAMP WITH TIME ZONE NULL,
    is_active   BOOLEAN DEFAULT FALSE,

    CONSTRAINT pk_user_account_id PRIMARY KEY (id),
    CONSTRAINT fk_user_account_acces_id FOREIGN KEY (acces_id) REFERENCES access (id)
);

-- +migrate StatementEnd