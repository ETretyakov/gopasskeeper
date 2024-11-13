BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS public.sec_cards
(
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    uid uuid NOT NULL,
    name varchar NULL,
    number varchar NOT NULL,
    mask varchar NOT NULL,
    month integer NULL,
    year integer NULL,
    cvc varchar NULL,
    pin varchar NULL,
    created_at timestamp WITHOUT time ZONE DEFAULT current_timestamp NOT NULL,
    CONSTRAINT sec_cards_pk PRIMARY KEY (id),
    CONSTRAINT fk_uid
      FOREIGN KEY(uid)
        REFERENCES public.usr_users(id)
);

COMMIT;