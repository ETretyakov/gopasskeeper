BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS public.syn_timestamps
(
    uid uuid NOT NULL,
    timestamp timestamp WITHOUT time ZONE NOT NULL,
    CONSTRAINT syn_timestamps_pk PRIMARY KEY (uid),
    CONSTRAINT fk_uid
      FOREIGN KEY(uid)
        REFERENCES public.usr_users(id)
);

COMMIT;