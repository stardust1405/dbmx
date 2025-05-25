-- Test SQL queries to demonstrate cursor-based query selection

SELECT * FROM users
WHERE age > 18
ORDER BY name

SELECT
    u.name,
    u.email,
    p.title as position
FROM users u
JOIN positions p ON u.position_id = p.id
WHERE u.active = true

UPDATE users
SET last_login = NOW()
WHERE id = 123

SELECT COUNT(*) as total_users,
       AVG(age) as average_age,
       MAX(created_at) as latest_signup
FROM users
WHERE active = true
GROUP BY department
