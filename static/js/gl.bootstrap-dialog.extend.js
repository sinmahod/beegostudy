(function ($) {
 /**
 * BootstrapDialog 封装类
 * author: gl
 * time: 2016/11/18 22:33:07
 * email: sinmahod@qq.com
 */
    window.BootFrame = function (e) {
	return {
		//弹出框
		alert: function (msg,fn,tle,iswarning) {
	        	BootstrapDialog.alert({
		            title: tle ? tle : '提示',
		            message: msg,
		            type: iswarning ? BootstrapDialog.TYPE_WARNING : BootstrapDialog.TYPE_PRIMARY, // <-- Default value is BootstrapDialog.TYPE_WARNING
		            closable: true, // <-- Default value is false
		            draggable: true, // <-- Default value is false
		            buttonLabel: '确定', // <-- Default value is 'OK',
		            callback: function(result) {
		                 	if (fn && fn instanceof Function) {
		               		fn();
	               		}
		            }
		        });
		},
		//选择框
		confirm: function (message,fn,falsefn){
			BootstrapDialog.confirm({
			            title: '请选择',
			            message: message,
			            type: BootstrapDialog.TYPE_WARNING, // <-- Default value is BootstrapDialog.TYPE_PRIMARY
			            closable: true, // <-- Default value is false
			            draggable: true, // <-- Default value is false
			            btnCancelLabel: '取消', // <-- Default value is 'Cancel',
			            btnOKLabel: '确定', // <-- Default value is 'OK',
			            btnOKClass: 'btn-warning', // <-- If you didn't specify it, dialog type will be used,
			            callback: function(result) {
					if(result) {
					   	if (fn && fn instanceof Function) {
			               			fn();
		               			}
					}else {
					    	if (falsefn && falsefn instanceof Function) {
			               			falsefn();
		               			}
					}
			             }
		        	});
		},
		//Dialog
		dialog: function(){
			var t,m;
			var b = {};
			var hideclose = false;
			var w = 600;
			var h = 150;
			var diaid;
			var id;
			var dobj;
			return{
				id : function(id){
					diaid = id;
				},
				title : function(title){
					t = title;
				},
				message :function(message){
					m = message;
				},
				url: function(u){
					m = $('<div></div>').load(u);
				},
				width:function(width){
					w = width / 5 * 5;
				},
				height:function(height){
					h = height / 5 * 5 ; //框头框尾130px
				},
				addButton: function(name,fn,css,keycode){
					var bf = this;  //获得dialog实体

					if(!css){
						css = 'btn-primary';
					}
					if (jQuery.isEmptyObject(b) ) {
						b = [{
							label:name,
							cssClass:css,
							hotkey:keycode?keycode:0,
							action:function(){
								if (fn && fn instanceof Function) {
						               		fn(bf,this);  //这里的this指的是按钮
					               		}
							}
						}]
					 }else{
						var bt = {
							label:name,
							cssClass:css,
							hotkey:keycode?keycode:0,
							action:function(){
								if (fn && fn instanceof Function) {
						               		fn(bf,this);  //这里的this指的是按钮
					               		}
							}
						}
						b.push(bt)
					}
				},
				hideClose: function(){
					hideclose = true;
				},
				show:function(){
					if (!hideclose){
						var bt = {
							label: '取消',
					                    	action: function(dialogItself){
					                        		dialogItself.close();
					                    	}
						}
						b.push(bt);
					}
				  	dobj = BootstrapDialog.show({
			  			   id: diaid,
					                title: t,
					                message: m,
					                cssClass: 'dialog-'+w+' dialog-h-'+h,
					                draggable: true,
					                buttons: b
				            });
				  	id = dobj.$modal[0].id;
				},
				getFormData: function(){
					var s = $("#"+id + " form");
					if (s.length == 0) {
						BootFrame.alert("未找到form标签");
						return;
					}
					if (s[0].length == 0 ){
						BootFrame.alert("form中不存在表单");
						return;
					}
					var a = s[0];
					var $map = new Map();    
					for (var i = 0 ; i < a.length; i ++){
  						$map.put(a[i].id , a[i].value);   
					}
					return $map;
				},
				close: function(){
					dobj.close();
				}
			}//<! return >
		}//<! dialog >
	}//<! return >
   }();

})(jQuery);