<?php $_v=&$this->vars; $this->push();$this->load("html_header.tpl");$this->assign($_v);$this->render();$this->pop();?>


<div class="container">

	<div class="page-head">
		<div>
			<h2><?php echo htmlspecialchars($_v['title'], ENT_QUOTES);?>
</h2>
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
<?php if(is_array($_v['output'])){?>

<?php if(!empty($_v['output']))foreach($_v['output'] as $_v['line']){?><div title="<?php echo $_v['line']['timestamp'];?>
" class="output-line <?php echo $_v['line']['level'];?>
 <?php echo $_v['line']['fields']['output'];?>
"><?php echo htmlspecialchars($_v['line']['message'], ENT_QUOTES);?>
</div>
<?php }else{?>
<div class="empty-cell">No output recorded.</div>
<?php }?>
<?php }else{?>
<div class="output-line"><?php echo htmlspecialchars($_v['output'], ENT_QUOTES);?>
</div>
<?php }?>
			</div>
		</div>
	</section>

</div>

<?php $this->push();$this->load("html_footer.tpl");$this->assign($_v);$this->render();$this->pop();?>

