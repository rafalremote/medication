-- Seed 1000 medications into the database
DO $$ 
BEGIN
    FOR i IN 1..1000 LOOP
        INSERT INTO medications (name, dosage, form) 
        VALUES (
            'Medication ' || i, 
            (10 + (i % 500)) || 'mg', 
            CASE WHEN i % 3 = 0 THEN 'Tablet' WHEN i % 3 = 1 THEN 'Capsule' ELSE 'Syrup' END
        );
    END LOOP;
END $$;
