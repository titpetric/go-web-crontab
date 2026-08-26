-- Fixture rows for the escaping suite. The job name and the description both
-- carry the double quote that delimits the attributes they are printed into,
-- plus a tag, so an unescaped value would close its attribute and add markup of
-- its own rather than merely look wrong.
--
-- Seeding happens after the service has started, because the schema is created
-- by the migrations the binary runs on boot.

DELETE FROM logs WHERE name = 'esc"job<i>';
DELETE FROM jobs WHERE name = 'esc"job<i>';

INSERT INTO jobs (name, description, created_at) VALUES
	('esc"job<i>', 'desc a"b<i>', '2026-08-26 12:00:00');

-- Logs are numbered newest-first by the dashboard, so the row below with the
-- highest stamp is row 1 and the one with the lowest is row 3. The newest is
-- also the run the index page reports as the last one, which is why it is the
-- failing one: it drives the fail badge and the red sparkline in one go.
INSERT INTO logs (name, stamp, duration, output, exit_code) VALUES
	('esc"job<i>', '2026-08-26 12:00:03', 1500000000, '[{"fields":{"job":"esc","output":"stderr"},"level":"error","timestamp":"2026-08-26T12:00:03Z\"<i>","message":"structured <script>alert(2)</script>"}]', 123),
	('esc"job<i>', '2026-08-26 12:00:02', 2500000000, 'plain <script>alert(1)</script>', 0),
	('esc"job<i>', '2026-08-26 12:00:01', 3500000000, 'plain output', 0);
