CREATE INDEX idx_sessions_user ON sessions(user_id,expires_at);
CREATE INDEX idx_tasks_project_status ON remediation_tasks(project_id,status,updated_at);
CREATE INDEX idx_samples_plot_time ON monitoring_samples(plot_id,collected_at);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type,entity_id,created_at);
CREATE INDEX idx_events_task ON task_events(task_id,created_at);
