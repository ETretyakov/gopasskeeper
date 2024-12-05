BEGIN TRANSACTION;

ALTER TABLE public.sec_accounts
    DROP column IF EXISTS meta;
ALTER TABLE public.sec_cards
    DROP column IF EXISTS meta;
ALTER TABLE public.sec_notes
    DROP column IF EXISTS meta;
ALTER TABLE public.sec_files
    DROP column IF EXISTS meta;

COMMIT;