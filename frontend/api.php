<?php

// @route GET /api/list

include("bootstrap.php");

$rows = $db->get_all("SELECT name, description FROM jobs WHERE deleted_at IS NULL OR deleted_at = '' ORDER BY name");

$body = array("jobs" => array());
foreach ($rows as $row) {
	$job = array(
		"name" => $row["name"],
		"description" => $row["description"],
		"active" => 1,
	);

	$last = $db->get("SELECT stamp, duration / 1000000000.0 as duration, exit_code as exitCode FROM logs WHERE name = ? ORDER BY stamp DESC LIMIT 1", $row["name"]);
	if ($last) {
		$job["lastRun"] = $last;
	}

	$job["history"] = $db->get_all("SELECT stamp, duration / 1000000000.0 as duration, exit_code as exitCode FROM logs WHERE name = ? ORDER BY stamp DESC LIMIT 20", $row["name"]);
	$body["jobs"][] = $job;
}

header("Content-Type: application/json");
echo json_encode($body);

$db->close();
