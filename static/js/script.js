//----- UI Iinitial -----//
$('input[name="date"]').val(getToday());
$('input[name="payer"]').hide()
$('#datetimepicker').datetimepicker({
    format: 'YYYY/MM/DD',
    defaultDate:new Date()
});
//----- socket.io -----//
var socket = io('', {transports: ['websocket']});//this domain
socket.on('added',function(data){
    alert("成功新增！")
    clear()
});
socket.on('err',function(data){
    alert("維修中...")
});
socket.on('list',function(data){
    var sheets = data.split("\\");
    var select = $('select')[0];
    var i=0
    /*
    sheets.forEach(sheet => {
        var vals = sheet.split("/")
        var option = new Option(vals[1], vals[0]);
        $(option).html(vals[1]);
        option.setAttribute("gid", vals[2])
        select.append(option);
        if(i==0){
            $('#link').text(vals[0]+":"+vals[1])
            $('#link').attr("href", "https://docs.google.com/spreadsheets/d/"+vals[0]+"/edit#gid="+vals[2])
        }
        i++
    });
    */
    sheets.forEach(sheet => {
        var vals = sheet.split("/")
        $("#sheetlist").append("<li><a href=\"javascript:selectSheet(\'"+sheet+"\')\">"+vals[1]+"</a></li>")
        if(i==0){
            selectSheet(sheet)
        }
        i++
    });
});
//----- ask list -----//
var data = {spreadsheetId:"1zvYlacc1ESyAcBoxuOyLlZ_Uiilz5MA8b21_p_NzWng",tableName:"List"}
socket.emit('requirelist', data)

//------------------------------------------------------------------------
$('select').on('change', function (){
    var s = $('select')[0]
    //var opt = s.options[s.selectedIndex]
    var opt = s.filter(":selected")
    $('#link').text(opt.value+":"+opt.text)
    $('#link').attr("href", "https://docs.google.com/spreadsheets/d/"+opt.value+"/edit#gid="+opt.attr("gid"))
});
//------------------------------------------------------------------------

$('input[name="state"]').on('change', function (){
    var x = this.value; // x gets the value attribute of changed checkbox
    if (this.checked){
        $('input[name="payer"]').show()
    }else {
        $('input[name="payer"]').hide()
    }
});
//------------------------------------------------------------------------
function selectSheet(sheet){
    var vals = sheet.split("/")
    $('#tabelName').attr("href", "https://docs.google.com/spreadsheets/d/"+vals[0]+"/edit#gid="+vals[2])
    $('#tabelName').text(vals[1])
    spreadsheetId=vals[0]
    tableName=vals[1]
}

var mode = false
function toggleMode(){
    mode = !mode
    clear()
    if(!mode){
        $("#optional").slideDown("fast")
        $("#submit").text("新增支出")
    }else{
        $("#optional").slideUp("fast")
        $("#submit").text("新增收入")
    }
}

function submit(){
    if(!mode){
        addOutcome()
    }else{
        addIncome()
    }
}
//------------------------------------------------------------------------
//********* global-var **********/
var tableName=""
var spreadsheetId=""
//*******************************/

function addOutcome(){
    if($("input[name='item']").val()==""){alert("名稱不能為空");return;}
    if($("input[name='amount']").val()==""){alert("金額不能為空");return;}
    var data = $('form').serializeObject();
    var s = $('select')[0]
    data.spreadsheetId = spreadsheetId
    data.tableName = tableName
    if(data.state=="on")data.state="Y"
    else data.state="N"
    if(data.receipt=="on")data.receipt="Y"
    else data.receipt="N"
    data.reimburse="N"
    data.outcome=data.amount
    data.income="0"
    //alert(JSON.stringify(data))
    socket.emit('add', data)
    /*
    socket.emit('add', {
      date  : "2018/11/27",
  		item  : "b-ball",
  	  payer : "Joseph",
  		state : "N",
  		reimburse : "Y",
  		income : "0",
  		outcome : "50",
    });
    */
}

function addIncome(){
    if($("input[name='item']").val()==""){alert("名稱不能為空");return;}
    if($("input[name='amount']").val()==""){alert("金額不能為空");return;}
    var data = $('form').serializeObject();
    var s = $('select')[0]
    data.spreadsheetId = spreadsheetId
    data.tableName = tableName
    data.state=""
    data.receipt=""
    data.reimburse=""
    data.income=data.amount
    data.outcome="0"
    //alert(JSON.stringify(data))
    socket.emit('add', data)
}

function clear(){
    $('input[name="date"]').val(getToday());
    $('input[name="item"]').val("");
    $('input[name="state"]').prop( "checked", false );;
    $('input[name="receipt"]').prop( "checked", false );;
    $('input[name="payer"]').val("");
    $('input[name="amount"]').val("");
    $('input[name="payer"]').hide()
}

$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
            o[this.name] = this.value || '';
        }
    });
    return o;
};

function getToday(){
  var today = new Date();
  var dd = today.getDate();
  var mm = today.getMonth() + 1; //January is 0!

  var yyyy = today.getFullYear();
  if (dd < 10) {
    dd = '0' + dd;
  }
  if (mm < 10) {
    mm = '0' + mm;
  }
  var today = yyyy + '-' + mm + '-' + dd;
  return today
}
