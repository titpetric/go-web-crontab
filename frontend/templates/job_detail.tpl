{load html_header.tpl}

<div class="container">
	<div class="row">
		<div class="col-8">
			<h2>{title}</h2>
		</div>
		<div class="col-4 tar">
			<a class="btn btn-primary" href="javascript:history.back()">back</a>
		</div>
	</div>

	<section class="contents">
	<code class="output">
		{if is_array($output)}
		{foreach $output as $line}
		<div title="{line.timestamp}" class="output-line {line.level} {line.fields.output}">{line.message|escape}</div>
		{/foreach}
		{else}
			<div class="output-line">{output|escape}</div>
		{/if}
	</code>
	</section>

	<a class="btn btn-primary" href="javascript:history.back()">back</a>

</div>

<style>
code {
	padding: 3px;
	display: block;
	overflow-x: scroll;
	background: #000;
	color: #eee;
	font-size: 12px;
}
code * {
	font-family: 'Roboto Mono';
	white-space: pre;
}
.output-line.stderr {
	color: #c00;
	font-weight: bold;
}
</style>

{load html_footer.tpl}