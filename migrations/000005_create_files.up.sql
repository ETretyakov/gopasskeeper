BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS public.sec_files
(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    uid uuid NOT NULL,
    name varchar NULL,
    created_at timestamp WITHOUT time ZONE DEFAULT current_timestamp NOT NULL,
    CONSTRAINT sec_files_pk PRIMARY KEY (id),
    CONSTRAINT fk_uid
      FOREIGN KEY(uid)
        REFERENCES public.usr_users(id)
);

COMMIT;