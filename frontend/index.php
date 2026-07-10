<?php

// @route GET /{jobName...}

include("bootstrap.php");

if (isset($_PATH["jobName"]) && $_PATH["jobName"] != "") {
	$jobName = $_PATH["jobName"];
	$pageNumber = 0;
	if (isset($_GET["pageNumber"])) {
		$pageNumber = (int)$_GET["pageNumber"];
	}
	if ($pageNumber < 0) {
		$pageNumber = 0;
	}

	$pageSize = 20;
	if (isset($_GET["pageSize"])) {
		$requested = (int)$_GET["pageSize"];
		if ($requested == 50) {
			$pageSize = 50;
		}
		if ($requested == 100) {
			$pageSize = 100;
		}
	}
	$offset = $pageNumber * $pageSize;

	$total = 0;
	$totalRow = $db->get("SELECT COUNT(*) AS total FROM logs WHERE name = ?", $jobName);
	if ($totalRow) {
		$total = $totalRow["total"] + 0;
	}

	$job = array(
		"jobName" => $jobName,
		"pageNumber" => $pageNumber,
		"pageSize" => $pageSize,
		"total" => $total,
		"firstRow" => $total > 0 ? $offset + 1 : 0,
		"lastRow" => ($offset + $pageSize) < $total ? $offset + $pageSize : $total,
		"hasPrev" => $pageNumber > 0,
		"hasNext" => ($offset + $pageSize) < $total,
		"link" => "/" . $jobName,
	);

	// Logs are numbered newest-first (row 1 = most recent) to stay consistent
	// with the row lookup used by detail.php.
	$logs = array();
	$rows = $db->get_all("SELECT stamp, duration / 1000000000.0 AS duration, exit_code AS exitCode FROM logs WHERE name = ? ORDER BY stamp DESC LIMIT ? OFFSET ?", $jobName, $pageSize, $offset);
	$rowNumber = $offset;
	foreach ($rows as $log) {
		$rowNumber++;
		$exitCode = $log["exitCode"] + 0;
		if ($exitCode != 0) {
			$exitCode = '<span class="badge badge-fail">Exit: ' . $exitCode . '</span>';
		} else {
			$exitCode = '<span class="badge badge-ok">OK</span>';
		}
		$logs[] = array(
			"id" => $rowNumber,
			"link" => $job["link"] . "/" . $rowNumber,
			"exitCode" => $exitCode,
			"duration" => format_duration($log["duration"]),
			"date" => date("Y/m/d H:i", strtotime($log["stamp"])),
		);
	}

	$dayRows = $db->get_all("SELECT strftime('%Y-%m-%d %H:00:00', stamp) AS stamp, MIN(duration) / 1000000000.0 AS durationMin, AVG(duration) / 1000000000.0 AS durationAvg, MAX(duration) / 1000000000.0 AS durationMax FROM logs WHERE name = ? AND stamp >= datetime('now', '-1 day') GROUP BY strftime('%Y-%m-%d %H', stamp) ORDER BY stamp", $jobName);
	$daily = array(
		"labels" => array(),
		"datasets" => array(
			array("label" => "Minimum duration", "backgroundColor" => "#c7d2fe", "borderRadius" => 4, "data" => array()),
			array("label" => "Average duration", "backgroundColor" => "#6366f1", "borderRadius" => 4, "data" => array()),
			array("label" => "Maximum duration", "backgroundColor" => "#3730a3", "borderRadius" => 4, "data" => array()),
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
			array("label" => "Minimum duration", "backgroundColor" => "#c7d2fe", "borderRadius" => 4, "data" => array()),
			array("label" => "Average duration", "backgroundColor" => "#6366f1", "borderRadius" => 4, "data" => array()),
			array("label" => "Maximum duration", "backgroundColor" => "#3730a3", "borderRadius" => 4, "data" => array()),
		),
	);
	foreach ($monthRows as $row) {
		$monthly["labels"][] = $row["stamp"];
		$monthly["datasets"][0]["data"][] = $row["durationMin"];
		$monthly["datasets"][1]["data"][] = $row["durationAvg"];
		$monthly["datasets"][2]["data"][] = $row["durationMax"];
	}

	$title = $jobName;

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

$summary = array(
	"total" => 0,
	"scheduled" => 0,
	"onetime" => 0,
	"healthy" => 0,
	"failing" => 0,
);
foreach ($body["jobs"] as $job) {
	$summary["total"]++;
	if (strpos($job["name"], "onetime/") !== false) {
		$summary["onetime"]++;
	} else {
		$summary["scheduled"]++;
	}
	if (!empty($job["lastRun"])) {
		if (($job["lastRun"]["exitCode"] + 0) == 0) {
			$summary["healthy"]++;
		} else {
			$summary["failing"]++;
		}
	}
}

$title = "webcron dashboard";

$tpl->load("main.tpl");
$tpl->assign(array("title" => $title, "body" => $body, "summary" => $summary));
$tpl->render();

$db->close();
