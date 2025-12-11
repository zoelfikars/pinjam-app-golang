-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE status_type AS ENUM ('pending', 'approved', 'rejected');
CREATE TABLE IF NOT EXISTS pengajuan_pinjaman (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    nasabah_id uuid NOT NULL REFERENCES users(id),
    nik VARCHAR(20) NOT NULL,
    nama_lengkap VARCHAR(255) NOT NULL,
    alamat TEXT NOT NULL,
    no_telepon VARCHAR(20) NOT NULL,
    jumlah_pinjaman BIGINT NOT NULL,
    status status_type NOT NULL DEFAULT 'pending',
    catatan_admin TEXT,
    inspected_by uuid REFERENCES users(id),
    inspected_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pengajuan_pinjaman;
DROP TYPE status_type;
-- +goose StatementEnd
