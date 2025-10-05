#!/bin/bash
set -e  # para parar caso algum arquivo .sql dê erro

DB_NAME="database"
DB_USER="admin"
DB_HOST="localhost"
DB_PORT="5432"
DB_PASS="admin"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SQL_DIR="$SCRIPT_DIR/sql"

export PGPASSWORD="$DB_PASS"

# Garante que a tabela de migrations existe
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
CREATE TABLE IF NOT EXISTS migrations (
    id SERIAL PRIMARY KEY,
    filename TEXT UNIQUE NOT NULL,
    executed_at TIMESTAMP DEFAULT NOW()
);
EOF

# Executa cada arquivo .sql em ordem
for file in $(ls "$SQL_DIR"/*.sql | sort -V); do
    filename=$(basename "$file")

    # Verifica se já foi executado
    ALREADY_RUN=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT 1 FROM migrations WHERE filename = '$filename' LIMIT 1;")

    if [[ -n "$ALREADY_RUN" && "$ALREADY_RUN" =~ 1 ]]; then
        echo "Pulando: $filename (já executado)"
    else
        echo "🚀 Executando: $filename"
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
        # Registra na tabela
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c \
            "INSERT INTO migrations (filename) VALUES ('$filename');"
    fi
done

echo "Todas as migrations foram processadas!"
