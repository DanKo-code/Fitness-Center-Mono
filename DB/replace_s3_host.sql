DO $$
    DECLARE
        new_host TEXT := 'https://localhost'; -- Введите новый хост
    BEGIN
        -- Обновление колонок photo в каждой из таблиц
        UPDATE abonement
        SET photo = regexp_replace(photo, '//[^/]+:', '//' || substring(new_host from 'https?://([^:/]+)') || ':', 'g')
        WHERE photo IS NOT NULL AND photo <> '';

        UPDATE coach
        SET photo = regexp_replace(photo, '//[^/]+:', '//' || substring(new_host from 'https?://([^:/]+)') || ':', 'g')
        WHERE photo IS NOT NULL AND photo <> '';

        UPDATE service
        SET photo = regexp_replace(photo, '//[^/]+:', '//' || substring(new_host from 'https?://([^:/]+)') || ':', 'g')
        WHERE photo IS NOT NULL AND photo <> '';

        UPDATE "user"
        SET photo = regexp_replace(photo, '//[^/]+:', '//' || substring(new_host from 'https?://([^:/]+)') || ':', 'g')
        WHERE photo IS NOT NULL AND photo <> '';
    END $$;


