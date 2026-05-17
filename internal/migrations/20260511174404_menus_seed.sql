-- +goose Up
-- +goose StatementBegin

-- Insert groups
INSERT INTO menu_groups (name, display_order) VALUES
('Overview', 1),
('Users & Access', 2),
('Content Library', 3),
('Operations', 4),
('Business', 5),
('Support & Comms', 6),
('System', 7);

-- Insert top-level items (parent_id = NULL) for each group
-- We'll use a CTE to map group names to IDs
WITH group_ids AS (
    SELECT id, name FROM menu_groups
)
INSERT INTO menu_items (group_id, parent_id, label, href, icon, permission, display_order) VALUES
-- Overview > Dashboard
((SELECT id FROM group_ids WHERE name='Overview'), NULL, 'Dashboard', '/dashboard', 'LayoutDashboard', NULL, 1),

-- Users & Access > Users (parent)
((SELECT id FROM group_ids WHERE name='Users & Access'), NULL, 'Users', NULL, 'Users', NULL, 1),
-- Users & Access children
((SELECT id FROM group_ids WHERE name='Users & Access'), (SELECT id FROM menu_items WHERE label='Users' AND parent_id IS NULL), 'Users', '/dashboard/users', 'Users', NULL, 1),
((SELECT id FROM group_ids WHERE name='Users & Access'), (SELECT id FROM menu_items WHERE label='Users' AND parent_id IS NULL), 'Access Control', '/dashboard/access-control', 'Users', NULL, 2),
((SELECT id FROM group_ids WHERE name='Users & Access'), (SELECT id FROM menu_items WHERE label='Users' AND parent_id IS NULL), 'Roles & Permissions', '/dashboard/roles', 'Users', NULL, 3),
((SELECT id FROM group_ids WHERE name='Users & Access'), (SELECT id FROM menu_items WHERE label='Users' AND parent_id IS NULL), 'Audit Logs', '/dashboard/audit-logs', 'Users', NULL, 4),

-- Content Library > Movies & Series (leaf)
((SELECT id FROM group_ids WHERE name='Content Library'), NULL, 'Movies & Series', '/dashboard/movies', 'Film', NULL, 1),
-- Content Library > Music (parent)
((SELECT id FROM group_ids WHERE name='Content Library'), NULL, 'Music', NULL, 'Music', NULL, 2),
-- Music children
((SELECT id FROM group_ids WHERE name='Content Library'), (SELECT id FROM menu_items WHERE label='Music' AND parent_id IS NULL), 'Artists', '/dashboard/music/artists', 'Music', NULL, 1),
((SELECT id FROM group_ids WHERE name='Content Library'), (SELECT id FROM menu_items WHERE label='Music' AND parent_id IS NULL), 'Albums', '/dashboard/music/albums', 'Music', NULL, 2),
((SELECT id FROM group_ids WHERE name='Content Library'), (SELECT id FROM menu_items WHERE label='Music' AND parent_id IS NULL), 'Tracks', '/dashboard/music/tracks', 'Music', NULL, 3),
-- Content Library > Media Assets (parent)
((SELECT id FROM group_ids WHERE name='Content Library'), NULL, 'Media Assets', NULL, 'FolderArchive', NULL, 3),
-- Media Assets children
((SELECT id FROM group_ids WHERE name='Content Library'), (SELECT id FROM menu_items WHERE label='Media Assets' AND parent_id IS NULL), 'Media Library', '/dashboard/media', 'FolderArchive', NULL, 1),
((SELECT id FROM group_ids WHERE name='Content Library'), (SELECT id FROM menu_items WHERE label='Media Assets' AND parent_id IS NULL), 'Uploads', '/dashboard/media/uploads', 'FolderArchive', NULL, 2),

-- Operations > Streaming (parent)
((SELECT id FROM group_ids WHERE name='Operations'), NULL, 'Streaming', NULL, 'MonitorPlay', NULL, 1),
-- Streaming children
((SELECT id FROM group_ids WHERE name='Operations'), (SELECT id FROM menu_items WHERE label='Streaming' AND parent_id IS NULL), 'Active Streams', '/dashboard/streams', 'MonitorPlay', NULL, 1),
((SELECT id FROM group_ids WHERE name='Operations'), (SELECT id FROM menu_items WHERE label='Streaming' AND parent_id IS NULL), 'Playback Errors', '/dashboard/playback/errors', 'MonitorPlay', NULL, 2),
((SELECT id FROM group_ids WHERE name='Operations'), (SELECT id FROM menu_items WHERE label='Streaming' AND parent_id IS NULL), 'DRM & Licenses', '/dashboard/drm', 'MonitorPlay', NULL, 3),
-- Operations > Moderation (parent)
((SELECT id FROM group_ids WHERE name='Operations'), NULL, 'Moderation', NULL, 'AlertTriangle', NULL, 2),
-- Moderation children
((SELECT id FROM group_ids WHERE name='Operations'), (SELECT id FROM menu_items WHERE label='Moderation' AND parent_id IS NULL), 'User Reports', '/dashboard/moderation/reports', 'AlertTriangle', NULL, 1),
((SELECT id FROM group_ids WHERE name='Operations'), (SELECT id FROM menu_items WHERE label='Moderation' AND parent_id IS NULL), 'DMCA & Copyright', '/dashboard/moderation/dmca', 'AlertTriangle', NULL, 2),
((SELECT id FROM group_ids WHERE name='Operations'), (SELECT id FROM menu_items WHERE label='Moderation' AND parent_id IS NULL), 'Region Restrictions', '/dashboard/moderation/region-blocks', 'AlertTriangle', NULL, 3),

-- Business > Monetization (parent)
((SELECT id FROM group_ids WHERE name='Business'), NULL, 'Monetization', NULL, 'CreditCard', NULL, 1),
-- Monetization children
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Monetization' AND parent_id IS NULL), 'Subscription Plans', '/dashboard/subscriptions/plans', 'CreditCard', NULL, 1),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Monetization' AND parent_id IS NULL), 'User Subscriptions', '/dashboard/subscriptions/users', 'CreditCard', NULL, 2),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Monetization' AND parent_id IS NULL), 'Payments', '/dashboard/payments', 'CreditCard', NULL, 3),
-- Business > Analytics (parent)
((SELECT id FROM group_ids WHERE name='Business'), NULL, 'Analytics', NULL, 'BarChart3', NULL, 2),
-- Analytics children
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Analytics' AND parent_id IS NULL), 'Overview', '/dashboard/analytics', 'BarChart3', NULL, 1),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Analytics' AND parent_id IS NULL), 'Content Performance', '/dashboard/analytics/content', 'BarChart3', NULL, 2),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Analytics' AND parent_id IS NULL), 'User Behavior', '/dashboard/analytics/users', 'BarChart3', NULL, 3),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Analytics' AND parent_id IS NULL), 'Reports', '/dashboard/reports', 'BarChart3', NULL, 4),
-- Business > Discovery (parent)
((SELECT id FROM group_ids WHERE name='Business'), NULL, 'Discovery', NULL, 'Sparkles', NULL, 3),
-- Discovery children
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Discovery' AND parent_id IS NULL), 'Recommendations', '/dashboard/recommendations', 'Sparkles', NULL, 1),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Discovery' AND parent_id IS NULL), 'Featured Content', '/dashboard/featured', 'Sparkles', NULL, 2),
((SELECT id FROM group_ids WHERE name='Business'), (SELECT id FROM menu_items WHERE label='Discovery' AND parent_id IS NULL), 'A/B Testing', '/dashboard/ab-tests', 'Sparkles', NULL, 3),

-- Support & Comms > Support (leaf)
((SELECT id FROM group_ids WHERE name='Support & Comms'), NULL, 'Support', '/dashboard/support', 'Headphones', NULL, 1),
-- Support & Comms > Communication (parent)
((SELECT id FROM group_ids WHERE name='Support & Comms'), NULL, 'Communication', NULL, 'Bell', NULL, 2),
-- Communication children
((SELECT id FROM group_ids WHERE name='Support & Comms'), (SELECT id FROM menu_items WHERE label='Communication' AND parent_id IS NULL), 'Notifications', '/dashboard/notifications', 'Bell', NULL, 1),
((SELECT id FROM group_ids WHERE name='Support & Comms'), (SELECT id FROM menu_items WHERE label='Communication' AND parent_id IS NULL), 'Campaigns', '/dashboard/campaigns', 'Bell', NULL, 2),

-- System > Infrastructure (parent)
((SELECT id FROM group_ids WHERE name='System'), NULL, 'Infrastructure', NULL, 'Server', NULL, 1),
-- Infrastructure children
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Infrastructure' AND parent_id IS NULL), 'System Logs', '/dashboard/logs', 'Server', NULL, 1),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Infrastructure' AND parent_id IS NULL), 'Monitoring', '/dashboard/monitoring', 'Server', NULL, 2),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Infrastructure' AND parent_id IS NULL), 'Backups', '/dashboard/backups', 'Server', NULL, 3),
-- System > Settings (parent)
((SELECT id FROM group_ids WHERE name='System'), NULL, 'Settings', NULL, 'Settings', NULL, 2),
-- Settings children
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Settings' AND parent_id IS NULL), 'General', '/dashboard/settings/general', 'Settings', NULL, 1),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Settings' AND parent_id IS NULL), 'Localization', '/dashboard/settings/localization', 'Settings', NULL, 2),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Settings' AND parent_id IS NULL), 'Feature Flags', '/dashboard/settings/features', 'Settings', NULL, 3),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Settings' AND parent_id IS NULL), 'Integrations', '/dashboard/settings/integrations', 'Settings', NULL, 4),
-- System > Legal (parent)
((SELECT id FROM group_ids WHERE name='System'), NULL, 'Legal', NULL, 'FileText', NULL, 3),
-- Legal children
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Legal' AND parent_id IS NULL), 'Terms of Service', '/dashboard/legal/terms', 'FileText', NULL, 1),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Legal' AND parent_id IS NULL), 'Privacy Policy', '/dashboard/legal/privacy', 'FileText', NULL, 2),
((SELECT id FROM group_ids WHERE name='System'), (SELECT id FROM menu_items WHERE label='Legal' AND parent_id IS NULL), 'Licenses', '/dashboard/legal/licenses', 'FileText', NULL, 3);

-- +goose StatementEnd