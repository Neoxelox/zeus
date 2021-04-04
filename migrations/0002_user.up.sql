CREATE TABLE "users" (
    "id"            VARCHAR(20) PRIMARY KEY,
    "name"          VARCHAR(100) NOT NULL,
    "username"      VARCHAR(100) UNIQUE NOT NULL,
    "age"           INTEGER NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "deleted_at"    TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX "users_name_idx" ON "users" USING gin ("name" gin_trgm_ops); -- CONCURRENTLY
CREATE INDEX "users_username_idx" ON "users" USING gin ("username" gin_trgm_ops); -- CONCURRENTLY
