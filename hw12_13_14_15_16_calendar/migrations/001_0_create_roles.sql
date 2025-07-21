DO $$
    BEGIN

        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'postgres_user') THEN
            CREATE USER postgres_user WITH PASSWORD 'postgres_password';
            RAISE NOTICE 'Role "postgres_user" created successfully';
        ELSE
            RAISE NOTICE 'Role "postgres_user" already exists';
        END IF;

    END
$$;


DO $$
    BEGIN

        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'notify_user') THEN
            CREATE USER notify_user WITH PASSWORD 'notify_password';
            RAISE NOTICE 'Role "notify_user" created successfully';
        ELSE
            RAISE NOTICE 'Role "notify_user" already exists';
        END IF;

    END
$$;
