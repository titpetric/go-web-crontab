<?php

include("vendor/autoload.php");

function route_job_name() {
	if (isset($_PATH["group"]) && $_PATH["group"] != "") {
		return $_PATH["group"] . "/" . $_PATH["jobName"];
	}

	return $_PATH["jobName"];
}

function format_duration($duration) {
	return sprintf("%.3fs", $duration);
}

function history($history) {
	$retval = array();
	for ($i = count($history) - 1; $i >= 0; $i--) {
		$run = $history[$i];
		$retval[] = array(
			"name" => isset($run["exitCode"]) ? "Exit: " . $run["exitCode"] : "OK",
			"date" => isset($run["stamp"]) ? $run["stamp"] : "",
			"value" => isset($run["duration"]) ? $run["duration"] : 0,
		);
	}

	return json_encode($retval);
}

$db = new Database("crontab");

$tpl = new MiniTPL\Template;

$tpl->set_paths('templates/');
$tpl->set_compile_location('cache/', false);
