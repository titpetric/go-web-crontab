<?php

// @route GET /{jobName...}

include("bootstrap.php");

if (isset($_PATH["jobName"]) && $_PATH["jobName"] != "") {
	$jobName = $_PATH["jobName"];
	$pageNumber = 0;
	if (isset($_GET["pageNumber"])) {
		$pageNumber = $_GET["pageNumber"] + 0;
	}
	$pageSize = 25;
	$offset = $pageNumber * $pageSize;

	$job = array(
		"jobName" => $jobName,
		"pageNumber" => $pageNumber,
		"pageSize" => $pageSize,
		"link" => "/" . $jobName,
	);

	$logs = array();
	$rows = $db->get_all("SELECT ID, stamp, duration, exitCode FROM (SELECT ROW_NUMBER() OVER (ORDER BY stamp DESC) AS ID, stamp, duration / 1000000000.0 AS duration, exit_code AS exitCode FROM logs WHERE name = ?) ORDER BY ID LIMIT ? OFFSET ?", $jobName, $pageSize, $offset);
	foreach ($rows as $log) {
		$exitCode = $log["exitCode"] + 0;
		if ($exitCode != 0) {
			$exitCode = '<span class="badge badge-danger">Exit: ' . $exitCode . '</span>';
		} else {
			$exitCode = '<span class="badge badge-success">OK</span>';
		}
		$logs[] = array(
			"id" => $log["ID"],
			"link" => $job["link"] . "/" . $log["ID"],
			"exitCode" => $exitCode,
			"duration" => format_duration($log["duration"]),
			"date" => date("Y/m/d H:i", strtotime($log["stamp"])),
		);
	}

	$dayRows = $db->get_all("SELECT strftime('%Y-%m-%d %H:00:00', stamp) AS stamp, MIN(duration) / 1000000000.0 AS durationMin, AVG(duration) / 1000000000.0 AS durationAvg, MAX(duration) / 1000000000.0 AS durationMax FROM logs WHERE name = ? AND stamp >= datetime('now', '-1 day') GROUP BY strftime('%Y-%m-%d %H', stamp) ORDER BY stamp", $jobName);
	$daily = array(
		"labels" => array(),
		"datasets" => array(
			array("label" => "Minimum duration", "backgroundColor" => "#caf270", "data" => array()),
			array("label" => "Average duration", "backgroundColor" => "#45c490", "data" => array()),
			array("label" => "Maximum duration", "backgroundColor" => "#008d93", "data" => array()),
		),
	);
	foreach ($dayRows as $row) {
		$x = explode(" ", $row["stamp"]);
		$daily["labels"][] = substr($x[1], 0, -3);
		$daily["datasets"][0]["data"][] = $row["durationMin"];
		$daily["datasets"][1]["data"][] = $row["durationAvg"];
		$daily["datasets"][2]["data"][] = $row["durationMax"];
	}

	$monthRows = $db->get_all("SELECT date(stamp) AS stamp, MIN(duration) / 1000000000.0 AS durationMin, AVG(duration) / 1000000000.0 AS durationAvg, MAX(duration) / 1000000000.0 AS durationMax FROM logs WHERE name = ? AND stamp >= datetime('now', '-30 days') GROUP BY date(stamp) ORDER BY stamp", $jobName);
	$monthly = array(
		"labels" => array(),
		"datasets" => array(
			array("label" => "Minimum duration", "backgroundColor" => "#caf270", "data" => array()),
			array("label" => "Average duration", "backgroundColor" => "#45c490", "data" => array()),
			array("label" => "Maximum duration", "backgroundColor" => "#008d93", "data" => array()),
		),
	);
	foreach ($monthRows as $row) {
		$monthly["labels"][] = $row["stamp"];
		$monthly["datasets"][0]["data"][] = $row["durationMin"];
		$monthly["datasets"][1]["data"][] = $row["durationAvg"];
		$monthly["datasets"][2]["data"][] = $row["durationMax"];
	}

	$title = $jobName . ", page " . $pageNumber;

	$tpl->load("job_stats.tpl");
	$tpl->assign(array("title" => $title, "daily" => $daily, "monthly" => $monthly, "job" => $job, "logs" => $logs));
	$tpl->render();

	$db->close();
	exit;
}

$rows = $db->get_all("SELECT name, description FROM jobs WHERE deleted_at IS NULL OR deleted_at = '' ORDER BY name");

$body = array("jobs" => array());

foreach ($rows as $row) {
	$job = array(
		"name" => $row["name"],
		"description" => $row["description"],
		"active" => 1,
	);

	$last = $db->get("SELECT stamp, duration / 1000000000.0 as duration, exit_code as exitCode FROM logs WHERE name = ? ORDER BY stamp DESC LIMIT 1", $job["name"]);
	if ($last) {
		$job["lastRun"] = $last;
	}

	$history = $db->get_all("SELECT stamp, duration / 1000000000.0 as duration, exit_code as exitCode FROM logs WHERE name = ? ORDER BY stamp DESC LIMIT 20", $job["name"]);
	$job["history"] = array();
	foreach ($history as $run) {
		$job["history"][] = $run;
	}

	$body["jobs"][] = $job;
}

$title = "CRONTAB jobs";

$tpl->load("main.tpl");
$tpl->assign(array("title" => $title, "body" => $body));
$tpl->render();

$db->close();
