<?php $_v=&$this->vars; $this->push();$this->load("html_header.tpl");$this->assign($_v);$this->render();$this->pop();?>


<div class="container">
	<div class="row">
		<div class="col-8">
			<h2><?php echo $_v['title'];?>
</h2>
		</div>
		<div class="col-4 tar">
			<a class="btn btn-primary" href="javascript:history.back()">back</a>
		</div>
	</div>

	<section class="contents">
	<code class="output">
		<?php if(is_array($_v['output'])){?>

		<?php if(!empty($_v['output']))foreach($_v['output'] as $_v['line']){?>
		<div title="<?php echo $_v['line']['timestamp'];?>
" class="output-line <?php echo $_v['line']['level'];?>
 <?php echo $_v['line']['fields']['output'];?>
"><?php echo htmlspecialchars($_v['line']['message'], ENT_QUOTES);?>
</div>
		<?php }?>
		<?php }else{?>

			<div class="output-line"><?php echo htmlspecialchars($_v['output'], ENT_QUOTES);?>
</div>
		<?php }?>
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

<?php $this->push();$this->load("html_footer.tpl");$this->assign($_v);$this->render();$this->pop();?>
