-- Seed the core configuration values
INSERT INTO system_configs (config_key, config_value) VALUES 
('current_session', '2026/2027'),
('current_semester', 'First'),
('registration_window', 'open');

-- Seed Clearance Offices dynamically
INSERT INTO clearance_offices (office_name) VALUES 
('Bursary'), ('Library'), ('Sports Department'), ('Faculty Dean Office');

-- Seed a sample Department
INSERT INTO departments (name, faculty) VALUES 
('Computer Science', 'Science and Technology');