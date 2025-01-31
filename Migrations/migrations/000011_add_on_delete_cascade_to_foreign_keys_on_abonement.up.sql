ALTER TABLE "order" DROP CONSTRAINT IF EXISTS order_abonement_id_fkey;
ALTER TABLE "order" DROP CONSTRAINT IF EXISTS order_user_id_fkey;

ALTER TABLE "order"
    ADD CONSTRAINT order_abonement_id_fkey
        FOREIGN KEY (abonement_id) REFERENCES abonement(Id) ON DELETE CASCADE;

ALTER TABLE "order"
    ADD CONSTRAINT order_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE;