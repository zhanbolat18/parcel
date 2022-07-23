-- +goose Up
-- +goose StatementBegin
INSERT INTO users(email, password_hash, status, role)
VALUES('admin@admin.com','$2a$13$JZk0FnS5tFr5gOLDvpOV9.6jzhD.6QiCpc267.y4TbVcR1/9KcN8W', 'active', 'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email='admin@admin.com' AND role='admin';
-- +goose StatementEnd
