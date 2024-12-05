BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS public.usr_users
(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    login varchar NOT NULL UNIQUE,
    pass_hash varchar NOT NULL,
    created_at timestamp WITHOUT time ZONE DEFAULT current_timestamp NOT NULL,
    CONSTRAINT usr_users_pk PRIMARY KEY (id),
	CONSTRAINT usr_users_unique UNIQUE (login)
);

COMMIT;