

$(window).resize(function(e) {
	$(".loginbox").height($(window).height());
    $(".main").height($(window).height() - $(".header").height() - $(".footer").height()-1);
	$(".left").height($(".main").height());
	$(".right").height($(".main").height());
	$(".rightdown").height($(".right").height()-27);
	$("#iframe").height($(window).height() - $(".header").height() - $(".footer").height()-29);
}).resize();



$(".menuup").click(function(){
  $(".headdown").slideUp("fast");
window.parent.sub_size();

});


$(function(){
    $("ul.toolbar li a").click(function(){
	 $("ul.toolbar li a").removeClass("selected")
	$(this).addClass("selected");
});
})


$(".nav li").click(function(){
  $(".headdown").slideDown("fast");
window.parent.add_size();
});


function sleep(delay) {
  var start = (new Date()).getTime();
  while ((new Date()).getTime() - start < delay) {
    continue;
  }
}

function alert_ajax1(url,formid){
	var data=getform1(formid);
	$.ajax({
			type:"POST",
			url:url,
			data:data,
			datatype:"json",
			success:function(ret){
				alert(ret);	
				close1();
				window.parent.reload1();
			}
	});
}


function loadn(left,right){
	window.parent.left.location.href=left
	window.parent.right.location.href=right
}
function reSizeList(){
	$("#list2").setGridHeight($(window).height()-135);
}
function getRows(){
	var idt=$("#list2").jqGrid('getGridParam','selrow')
	if(idt== null){
		alert("未选中行")
		return
	}
	var data=$("#list2").jqGrid("getRowData",idt)		
	return data
}
function getMultiRows(opt){
	var idt=$("#list2").jqGrid('getGridParam','selarrrow')
	if(idt== null){
		alert("未选中行")
		return
	}
	
	var i=0
	if(opt=="" || opt==undefined){
		var data3=[]
		$(idt).each(function(index,id){
			var data1=$("#list2").jqGrid("getRowData",id)
				data3[i]=data1
				i+=1
		})
		return data3
	}else{
		var data2=[]
		var data={}
		$(idt).each(function(index,id){
			var data1=$("#list2").jqGrid("getRowData",id)
				data2[i]=data1[opt]				
				i+=1
		})
		data[opt]=data2	
		return data
	}
}

function getMultiRows1(opt){
	var idt=$("#list2").jqGrid('getGridParam','selarrrow')
	if(idt== null){
		alert("未选中行")
		return
	}
	return idt
}


function alert_ajax_data(url,data){
	$.ajax({
			type:"POST",
			url:url,
			data:data,
			datatype:"json",
			success:function(ret){
				alert(ret);	
			}
	});
}

function alert_ajax_json(url,data){
	$.ajax({
			type:"POST",
			url:url,
			data:data,
			datatype:"json",
			contentType:"application/json;charset=UTF-8",
			success:function(ret){
				alert(ret);	
			}
	});
}

function alert_ajax_reload(url,data){
	$.ajax({
			type:"POST",
			url:url,
			data:data,
			datatype:"json",
			success:function(ret){
				alert(ret);	
				window.location.reload();
			}
	});
}

function getform1(formid){
	var d = {};
    	var t = $("#"+formid).serializeArray();
	$.each(t, function() {
      		d[this.name] = this.value;
	});
	return d;
}


function form(height,width,urlf){
	var index1=layer.open({
		title:'编辑',
		type: 2,
		content:urlf,
		area: [height+'px', width+'px'],
		fixed: false, //不固定
		maxmin: true,
		btnAlign: 'c',
		id:"pop"
	});
	index=index1;

}
function form1(height,width,data){

	layer.open({
	  title:'编辑用户',
	  type: 2,
	  //content: html,
	content:"add/pop",
	  area: ['600px', '450px'],
	  fixed: false, //不固定
	  maxmin: true,
	  btnAlign: 'c'
	
  });




}


function form2(height,width,urlf){
	var index1=layer.open({
		title:'编辑',
		type: 2,
		content:urlf,
		area: [height+'px', width+'px'],
		fixed: false, //不固定
		maxmin: true,
		scrollbar:true,
		resize:true,
		btnAlign: 'c',
		id:"pop"
	});
	index=index1;

}



$('.simpletable tbody tr:even').addClass('even');




