-- +goose Up
CREATE TABLE IF NOT EXISTS certificate_applications
(
    id SERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    application_status VARCHAR(10) NOT NULL CHECK(application_status IN('Pending', 'Rejected', 'Prepare', 'Done', 'Cancelled')) DEFAULT 'Pending',
    certificate_type VARCHAR(32) NOT NULL CHECK(certificate_type IN ('StudyPeriod','Recommendation','Academic','Common')) DEFAULT 'Common',
    obtain_method VARCHAR(16) NOT NULL CHECK(obtain_method IN('Paper', 'Electronic')) DEFAULT 'Electronic',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS certificates 
(
    id SERIAL PRIMARY KEY,
    receiver_id BIGINT DEFAULT NULL,
    author_id BIGINT NOT NULL,
    file_name TEXT DEFAULT NULL,
    storage_url TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cert_applications_student_id ON certificate_applications(student_id);
CREATE INDEX idx_cert_applications_status ON certificate_applications(application_status);
CREATE INDEX idx_cert_applications_created_at ON certificate_applications(created_at);

CREATE INDEX idx_certificates_receiver_id ON certificates(receiver_id);
CREATE INDEX idx_certificates_author_id ON certificates(author_id);
CREATE INDEX idx_certificates_created_at ON certificates(created_at);

-- +goose Down
DROP TABLE IF EXISTS certificate_applications;
DROP TABLE IF EXISTS certificates;
