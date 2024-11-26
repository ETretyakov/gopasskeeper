BEGIN TRANSACTION;

ALTER TABLE public.sec_accounts
    ADD column IF NOT EXISTS meta VARCHAR DEFAULT '';
ALTER TABLE public.sec_cards
    ADD column IF NOT EXISTS meta VARCHAR DEFAULT '';
ALTER TABLE public.sec_notes
    ADD column IF NOT EXISTS meta VARCHAR DEFAULT '';
ALTER TABLE public.sec_files
    ADD column IF NOT EXISTS meta VARCHAR DEFAULT '';

COMMIT;