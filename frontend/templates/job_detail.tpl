{load html_header.tpl}

<div class="container">

	<div class="page-head">
		<div>
			<h2>{title|escape}</h2>
			<div class="subtitle">Captured job output</div>
		</div>
		<div class="page-actions">
			<a class="btn btn-secondary" href="javascript:history.back()">&larr; Back</a>
			<a class="btn btn-secondary" href="/">Dashboard</a>
		</div>
	</div>

	<section class="panel">
		<div class="panel__head">
			<h2 class="panel__title">Output</h2>
		</div>
		<div class="panel__body--pad">
			<div class="log-output">
{if is_array($output)}
{foreach $output as $line}<div title="{line.timestamp}" class="output-line {line.level} {line.fields.output}">{line.message|escape}</div>
{else}<div class="empty-cell">No output recorded.</div>
{/foreach}
{else}<div class="output-line">{output|escape}</div>
{/if}
			</div>
		</div>
	</section>

</div>

{load html_footer.tpl}
